<template>
  <input class="upload-select" ref="selectRef" type="file" multiple @change="fileSelect" />
  <ElDropdown :teleported="false" @command="handleUpload">
    <el-button type="primary" @click="handleUpload('file')">
      <ElIcon :size="18">
        <i-mynaui:upload />
      </ElIcon>
      <span>上传</span>
    </el-button>

    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="file">上传文件</el-dropdown-item>
        <el-dropdown-item command="directory">上传文件夹</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </ElDropdown>

  <!-- 传输列表 -->
  <div class="upload-panel">
    <el-button plain @click="togglePanel">
      <el-icon :size="16"><Sort /></el-icon>
      <span>传输列表</span>
    </el-button>

    <!-- panel -->
    <el-card v-show="showPanel" class="file-panel">
      <span class="select-text">正在上传（{{ resolvedCount }}/{{ uploadFileList.length }}）</span>
      <el-table :data="uploadFileList" max-height="25rem">
        <el-table-column prop="name" width="150" show-overflow-tooltip label="文件名">
          <template #default="{ row }">
            <el-icon style="vertical-align: -10%">
              <i-mynaui:file />
            </el-icon>
            {{ row.name }}
          </template>
        </el-table-column>
        <el-table-column label="当前状态" width="200">
          <template #default="{ row }">
            <template v-if="[0, 1].includes(row.status)">
              <el-icon class="is-loading">
                <Loading />
              </el-icon>
              {{ row.status === 0 ? '等待处理' : '正在解析中' }}
            </template>
            <el-progress v-else-if="[2, 3].includes(row.status)" :percentage="row.percent" />
            <template v-else-if="row.status === 4">
              <el-icon style="vertical-align: middle">
                <i-mynaui:check-circle />
              </el-icon>
              上传成功
            </template>
            <template v-else-if="row.status === 5">
              <el-icon style="vertical-align: middle">
                <i-mynaui:daze-circle />
              </el-icon>
              上传失败!
            </template>
          </template>
        </el-table-column>
        <el-table-column width="120" prop="size" label="大小" />
        <el-table-column width="200" label="操作">
          <template #default="{ row }">
            <el-button
              v-if="[2, 3].includes(row.status)"
              type="primary"
              size="small"
              @click="toggleStatus(row)">
              {{ row.status === 2 ? '暂停' : '继续' }}
            </el-button>
            <el-button
              v-if="![0, 1].includes(row.status)"
              type="primary"
              size="small"
              @click="cancel(row.fileId)">
              {{ row.status === 4 ? '清除' : '取消' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>

  <div class="panel-cover" v-show="showPanel" @click.self="showPanel = false"></div>
</template>

<script setup lang="ts">
  import { ref, computed, shallowReactive } from 'vue';
  import { Sort, Loading } from '@element-plus/icons-vue';
  import { useRoute } from 'vue-router';
  import { formatFileSize } from '@/utils/filestatus';
  import { verify, uploadChunk, mergeChunks } from '@/api/file/fileList';
  import type { UploadFileItem, workerResult } from './types';
  import { cutFile } from './cutFile';

  defineOptions({ name: 'FileUpload' });
  const props = defineProps({ onSuccessed: { type: Function, default: () => {} } });
  const route = useRoute();
  const resolvedCount = computed(
    () => uploadFileList.value.filter((item) => item.status === 4).length
  ); // 已上传的文件数量

  // 上传文件or文件夹
  const selectRef = ref<HTMLInputElement>(null);
  const handleUpload = (cmd: string, e?: MouseEvent) => {
    const uploadInp = selectRef.value;
    if (cmd === 'file') {
      uploadInp.removeAttribute('webkitdirectory');
      uploadInp.removeAttribute('mozdirectory');
      uploadInp.removeAttribute('odirectory');
    } else {
      uploadInp.setAttribute('webkitdirectory', '');
      uploadInp.setAttribute('mozdirectory', '');
      uploadInp.setAttribute('odirectory', '');
    }
    uploadInp.click();
  };
  const uploadFileList = ref<UploadFileItem[]>([]);
  const uploadingData = new Map(); // 用于记录正在上传的文件的数据
  const fileSelect = (e: InputEvent) => {
    const files = Array.from((e.target as HTMLInputElement).files!);
    if (!files || files.length === 0) return;

    for (const file of files) {
      const { name, size } = file;
      const fileSize = formatFileSize(false, size);
      const fileItem = shallowReactive({ name, size: fileSize, status: 0, percent: 0, fileId: '' });
      uploadFileList.value.push(fileItem);
      fileItem.status = 1;
      showPanel.value = true;
      // 1. 分片并计算hash
      cutFile(file).then((result: workerResult) => {
        const { fileName, fileId, chunks } = result;
        fileItem.fileId = fileId;
        // 2. 确定需要上传的分片
        verify(fileName, fileId, chunks).then(({ code, shouldUpload }) => {
          const uploadList = chunks.filter((chunk) => shouldUpload.includes(chunk.chunkId));
          // 初始化文件上传信息
          const uploadingItem = {
            totalChunks: chunks.length,
            uploadedChunks: chunks.length - uploadList.length,
            chunks: structuredClone(uploadList),
          };
          uploadingData.set(fileId, uploadingItem);
          setProgress(fileId);
          startUploadTasks(fileId, uploadList);
        });
      });
    }
    selectRef.value.value = null;
  };

  // 上传列表的显示
  const showPanel = ref(false);
  const togglePanel = () => (showPanel.value = !showPanel.value);

  const MAX_CONCURRENT_UPLOADS = 6; // 最大并发上传数
  let fileUploadQueues = [];
  let activeUploads = 0;
  function startUploadTasks(fileId, chunks) {
    const fileInfo = uploadFileList.value.find((file) => file.fileId === fileId);
    fileUploadQueues.push({ fileId, fileInfo, chunks, status: 'uploading' }); // 保存该文件的分片队列
    if (chunks.length == 0 ) {
      // 合并分片
      const path = route.query.path || '';
      doMergeChunks(fileId, path,fileInfo);
      return;
    }
    processUploadQueue();
  }
  function processUploadQueue() {
    let uploadQueue = fileUploadQueues.filter((item) => item.status === 'uploading');
    while (uploadQueue.length > 0 && activeUploads < MAX_CONCURRENT_UPLOADS) {
      uploadQueue.forEach(({ fileId, fileInfo, chunks }) => {
        if (activeUploads >= MAX_CONCURRENT_UPLOADS) return;
        if (chunks.length > 0) {
          const chunk = chunks.shift();
          fileInfo.status = 2;
          activeUploads++;
          uploadChunk(fileId, chunk)
            .then(() => {
              activeUploads--;
              // 计算上传进度
              const data = uploadingData.get(fileId);
              const { uploadedChunks, totalChunks, chunks } = data;
              uploadingData.set(fileId, { ...data, uploadedChunks: uploadedChunks + 1 });
              setProgress(fileId);
              if (uploadedChunks + 1 === totalChunks) {
                  // 合并分片
                  const path = route.query.path || '';
                  doMergeChunks(fileId, path,fileInfo);
              }
              processUploadQueue();
            })
            .catch((error) => {
              activeUploads--;
              processUploadQueue();
            });
        } else {
          // 移除已完成上传的文件队列
          fileUploadQueues = fileUploadQueues.filter((item) => item.fileId !== fileId);
          uploadQueue = fileUploadQueues.filter((item) => item.status === 'uploading');
        }
      });
    }
  }

  function doMergeChunks(fileId, path,fileInfo) {
    mergeChunks(fileId, path).then((needs) => {
      if (needs.length !== 0) {
        // 更新进度条
        const data = uploadingData.get(fileId);
        uploadingData.set(fileId, {
          ...data,
          uploadedChunks: totalChunks - needs.length,
        });
        setProgress(fileId);
        // 重新上传需要上传的分片
        const uploadList = chunks.filter((chunk) => needs.includes(chunk.chunkId));
        fileUploadQueues.push({
          fileId,
          fileInfo,
          chunks: uploadList,
          status: 'uploading',
        });
        processUploadQueue();
      } else {
        clearFileInfo(fileId);
        fileInfo.status = 4;
        props.onSuccessed();
        if (resolvedCount.value === uploadFileList.value.length) {
          showPanel.value = false;
        }
      }
    });
  }

  const clearFileInfo = (fileId) => {
    uploadingData.delete(fileId);
  };
  const setProgress = (fileId) => {
    const file = uploadFileList.value.find((file) => file.fileId === fileId);
    const { totalChunks, uploadedChunks } = uploadingData.get(fileId);
    file.percent = Number(((uploadedChunks / totalChunks) * 100).toFixed(2));
  };
  const toggleStatus = (row) => {
    const idx = fileUploadQueues.findIndex((item) => item.fileId === row.fileId);
    if (fileUploadQueues[idx].status === 'paused') {
      row.status = 2;
      fileUploadQueues[idx].status = 'uploading';
      processUploadQueue();
      return;
    }
    row.status = 3;
    fileUploadQueues[idx].status = 'paused';
  };
  const cancel = (fileId) => {
    const idx = fileUploadQueues.findIndex((item) => item.fileId === fileId);
    if (idx != -1) fileUploadQueues.splice(idx, 1);
    const fileInfoIdx = uploadFileList.value.findIndex((file) => file.fileId === fileId);
    if (fileInfoIdx != -1) uploadFileList.value.splice(fileInfoIdx, 1);
    clearFileInfo(fileId);
    processUploadQueue();
  };
</script>

<style scoped lang="scss">
  .upload-select {
    display: none;
  }
  .upload-panel {
    position: relative;
    .file-panel {
      position: absolute;
      right: 0;
      top: 2.3rem;
      z-index: 999;
      text-align: left;
      // max-height: ;
      .select-text {
        font-size: 1rem;
        font-weight: 600;
        &::after {
          content: '';
          display: block;
          height: 1px;
          background-color: var(--ep-border-color-lighter);
          margin: 0.3rem 0;
        }
      }
    }
  }
  .panel-cover {
    position: fixed;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
    z-index: 998;
  }
</style>
