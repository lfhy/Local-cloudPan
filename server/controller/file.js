const config = require('../config');
const multer = require('@koa/multer'); // 文件上传
const archiver = require('archiver'); // 文件压缩
const sharp = require('sharp'); // 图片压缩
const send = require('koa-send');
const fs = require('fs-extra');
const path = require('path');
const Router = require('@koa/router');
const { sortByName, sortBySize, sortByModified, getNewFileName } = require('./utils');
const upload = multer({
  storage: multer.diskStorage({
    destination: function (req, _, cb) {
      const uploadPath = path.join(config.global.publicPath, req.body.path);
      req.uploadPath = uploadPath;
      cb(null, uploadPath);
    },
    filename: async function (req, file, cb) {
      // 解决中文乱码问题
      const filename = Buffer.from(file.originalname, 'latin1').toString('utf-8');
      const fullPath = path.join(req.uploadPath, filename);
      const exist = await fs.pathExists(fullPath);
      if (exist) {
        const curPath = await getNewFileName(fullPath);
        cb(null, path.basename(curPath));
      } else {
        cb(null, filename);
      }
    },
  }),
  limits: undefined, // 不设置文件大小限制
});

// 文件上传
const router = new Router();
router.post('/upload', upload.single(config.single.fieldName), async (ctx) => {
  ctx.status = 200;
  ctx.body = {
    code: 200,
    msg: '文件上传成功',
  };
});

// 文件下载
router.get('/checkFileList', async (ctx) => {
  try {
    let { 'filenameList[]': fList, path: reqPath } = ctx.request.query;
    fList = Array.isArray(fList) ? [...fList] : [fList];
    reqPath = path.join(config.global.publicPath, reqPath);
    if (!fList || fList.length === 0) {
      ctx.body = { code: 201, msg: 'Filename list is required' };
      return;
    }
    const filePathList = fList.map((filename) => path.join(reqPath, filename));
    const existenceChecks = filePathList.map((filePath) => fs.pathExists(filePath));
    const existenceResults = await Promise.all(existenceChecks);
    const missingFiles = fList.filter((_, index) => !existenceResults[index]);
    if (missingFiles.length > 0) {
      ctx.body = { code: 201, msg: `${missingFiles.join(', ')} 文件不存在！本次操作无效！` };
      return;
    } else {
      ctx.body = { code: 200, msg: 'ok' };
    }
  } catch (error) {
    ctx.status = 500;
  }
});
router.get('/download', async (ctx) => {
  try {
    let { filenameList: fList, path: reqPath } = ctx.request.query;
    fList = fList.slice(1, -1).split(',');
    reqPath = path.join(config.global.publicPath, reqPath);
    ctx.res.onerror = (err) => {
      console.log('Error while sending file:', err);
    };
    const isDir = (filename) => fs.statSync(path.join(reqPath, filename)).isDirectory();
    // 是否需要进行压缩下载
    const shouldCompress = fList.length > 1 || isDir(fList[0]);
    if (shouldCompress) {
      ctx.status = 200;
      ctx.attachment(`${fList[0]} 等文件.zip`);

      const archive = archiver('zip', {
        zlib: { level: 9 }, // 设置压缩级别
      });
      archive.on('error', function (err) {
        console.log(err);
        ctx.throw(500, 'Internal Server Error');
      });

      // 从Koa的响应对象的res中获取可写流，将其作为pipe的目标
      archive.pipe(ctx.res);

      // 添加文件或目录到压缩包
      try {
        for (const filename of fList) {
          const filepath = path.join(reqPath, filename);
          if (isDir(filename)) {
            archive.directory(filepath, { name: filename });
          } else {
            // 可以在这里添加文件大小检查
            archive.file(filepath, { name: filename });
          }
        }
        // 完成压缩文件的添加
        await archive.finalize();
      } catch (error) {
        console.error('Error during archiving:', error);
        ctx.throw(500, 'Internal Server Error');
      }
    }
    // 执行单个文件的下载
    else {
      ctx.attachment(fList[0]);
      ctx.length = fs.statSync(path.join(reqPath, fList[0])).size;
      await send(ctx, fList[0], { root: reqPath });
    }
  } catch (err) {
    console.log('Error while sending file:', err);
    ctx.status = 500;
  }
});

// 文件读取
const picType = ['jpeg', 'jpg', 'png', 'svg', 'gif', 'webp'];
router.get('/fileList', async (ctx) => {
  try {
    const { path: reqPath, sortMode } = ctx.request.query;
    const reqFilePath = path.join(config.global.publicPath, reqPath);
    // 读取路径下所有文件
    const files = await fs.readdir(reqFilePath, {
      withFileTypes: true,
    });
    const fileList = [];
    // 获取具体信息
    for (const file of files) {
      const name = file.name;
      const isDir = file.isDirectory() ? true : false;
      const ext = path.extname(name).substring(1).toLowerCase();
      let thumbnailPath;
      // 支持压缩的图片类型
      const reqpath = reqPath === '' ? '' : reqPath + '/';
      if (picType.includes(ext) && !isDir) {
        thumbnailPath = `api/thumbnail/${reqpath}${name}`;
      }
      const filePath = `api/${reqpath}${name}`;

      try {
        const { ino, size, mtimeMs } = await fs.stat(path.join(reqFilePath, name));
        fileList.push({
          id: ino,
          name,
          isDir,
          ext,
          size,
          modified: mtimeMs,
          filePath,
          thumbnailPath,
        });
      } catch (e) {
        // 文件无权限或错误的文件路径
      }
    }
    // 排序
    switch (sortMode) {
      case 'name':
        fileList.sort((a, b) => sortByName(a, b));
        break;
      case 'size':
        fileList.sort((a, b) => sortBySize(a, b));
        break;
      case 'modified':
        fileList.sort((a, b) => sortByModified(a, b));
        break;
    }
    ctx.res.status = 200;
    ctx.body = { code: 200, data: { fileList } };
  } catch (error) {
    console.log(error);
    ctx.response.status = 500;
  }
});

// 文件删除
router.delete('/delete', async (ctx) => {
  try {
    let { 'filenameList[]': fList, path: reqPath } = ctx.request.query;
    reqPath = path.join(config.global.publicPath, reqPath);
    fList = Array.isArray(fList) ? [...fList] : [fList];
    if (!fList || fList.length === 0) {
      ctx.body = { code: 201, msg: '请先选择文件！' };
      return;
    }
    const filePathList = fList.map((filename) => path.join(reqPath, filename));
    const existenceChecks = filePathList.map((filePath) => fs.pathExists(filePath));
    const existenceResults = await Promise.all(existenceChecks);

    const missingFiles = fList.filter((_, index) => !existenceResults[index]);
    if (missingFiles.length > 0) {
      ctx.body = { code: 201, msg: `${missingFiles.join(', ')} 文件不存在！本次操作无效！` };
      return;
    }
    // 执行删除操作
    await Promise.all(filePathList.map((filePath) => fs.remove(filePath)));
    ctx.res.status = 200;
    ctx.body = { code: 200, msg: '文件删除成功' };
  } catch (error) {
    console.log(error);
    ctx.response.status = 500;
  }
});

// 文件重命名
router.post('/rename', async (ctx) => {
  try {
    const { oldName, newName, path: reqPath } = ctx.request.body;
    const oldPath = path.join(config.global.publicPath, reqPath, oldName);
    const newPath = path.join(config.global.publicPath, reqPath, newName);
    const isExist = await fs.pathExists(oldPath);
    if (!isExist) {
      ctx.body = { code: 201, msg: `${oldName} 文件不存在！本次操作无效！` };
      return;
    }
    await fs.rename(oldPath, newPath);
    ctx.res.status = 200;
    ctx.body = { code: 200, msg: '文件重命名成功' };
  } catch (error) {
    console.log(error);
    ctx.response.status = 500;
  }
});

// 文件移动或复制
router.post('/fileMoveOrCopy', async (ctx) => {
  try {
    ctx.res.status = 200;
    const { fileList, destination, path: reqPath, dtype } = ctx.request.body;
    const operate = dtype === 'move' ? fs.move : fs.copy;
    const operateText = dtype === 'move' ? '移动' : '复制';

    // 错误检查
    for (const file of fileList) {
      const oldPath = path.join(config.global.publicPath, reqPath, file);
      const newPath = path.join(config.global.publicPath, destination, file);

      // 检查新路径是否包含旧路径
      if (newPath.includes(oldPath)) {
        ctx.body = { code: 201, msg: `不能${operateText}到自身或子目录!` };
        return;
      }

      // 检查旧文件是否存在
      const isExist = await fs.pathExists(oldPath);
      if (!isExist) {
        ctx.body = { code: 201, msg: `${file} 文件不存在！本次操作无效！` };
        return;
      }
    }

    // 执行文件移动或复制操作
    for (const file of fileList) {
      const oldPath = path.join(config.global.publicPath, reqPath, file);
      const newPath = path.join(config.global.publicPath, destination, file);
      await operate(oldPath, newPath, { overwrite: true });
    }
    ctx.body = { code: 200, msg: `文件${operateText}成功` };
  } catch (error) {
    console.log(error);
    ctx.response.status = 500;
  }
});

// 文件夹创建
router.post('/createDir', async (ctx) => {
  try {
    const { dirName, path: reqPath } = ctx.request.body;
    const dirPath = path.join(config.global.publicPath, reqPath, dirName);
    await fs.mkdir(dirPath, { recursive: true });
    ctx.res.status = 200;
    ctx.body = { code: 200, msg: '文件夹创建成功' };
  } catch (error) {
    console.log(error);
    ctx.response.status = 500;
  }
});

// 图片文件压缩
router.get('/thumbnail/:path*', async (ctx) => {
  const { path: reqPath } = ctx.params;
  sharp.cache(false);
  const imagePath = path.join(config.global.publicPath, reqPath);
  try {
    const image = sharp(imagePath, { animated: true });
    // 读取原始图片的元数据
    const { format, width } = await image.metadata();
    const processImage = image.rotate().resize(Math.round(width / 6));
    if (format === 'gif') {
      ctx.body = await processImage.gif().toBuffer();
    } else {
      ctx.body = await processImage.jpeg({ quality: 30 }).toBuffer();
    }
    ctx.type = `image/${format}`;
    ctx.set('Cache-Control', 'max-age=86400');
  } catch (error) {
    console.error('处理图片失败:', error);
    ctx.status = 500;
    ctx.body = '图片加载失败';
  }
});

module.exports = router;
