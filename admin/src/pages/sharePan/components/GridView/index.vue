<template>
  <div
    v-for="(file, index) in fileList"
    :key="file.id"
    class="file-item"
    :class="{ checked: checkedList[file.id] }"
    @contextmenu="handleContextMenu($event, file.id)"
    @click="fileChoosed(file.id)">
    <input class="file-checkbox" type="checkbox" v-model="checkedList[file.id]" name="selected" />
    <div @dblclick="fileClick(index)">
      <!-- 图片展示 -->
      <template v-if="file.fileType === 'picture'">
        <el-image
          :src="file.thumbnailPath || file.filePath"
          :lazy="true"
          :preview-src-list="[file.filePath]"
          class="picture"
          fit="cover"
          preview-teleported
          @click.stop="() => {}" />
      </template>
      <img v-else :src="getAssetsFile(file)" class="picture" />

      <!-- 文件名称 -->
      <template v-if="file.isRename">
        <input
          class="rename-ipt"
          type="text"
          name="filename"
          ref="renameIpt"
          v-focus
          :value="file.name"
          @click.stop="() => {}" />
        <ElButton
          type="primary"
          size="small"
          class="rename-btn"
          @click.stop="finishRenameOrCreate(file)">
          <ElIcon size="12"><Select /></ElIcon>
        </ElButton>
        <ElButton type="primary" size="small" class="rename-btn" @click.stop="cancelRename(file)">
          <ElIcon size="12"><CloseBold /></ElIcon>
        </ElButton>
      </template>
      <p v-else class="filename">{{ file.name }}</p>

      <!-- 文件描述 -->
      <div class="describe">
        <p v-if="file.isDir">{{ file.modified }}</p>
        <p v-else>{{ file.size }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, watch } from 'vue';
  import { getAssetsFile } from '@/utils/tool';
  import { Select, CloseBold } from '@element-plus/icons-vue';
  import type { formatFile } from '@/api/file/types';
  import { useRoute } from 'vue-router';

  defineOptions({ name: 'GridView' });
  const props = defineProps<{ fileList: formatFile[] }>();
  const emit = defineEmits(['fileClick', 'fileChoosed', 'finish', 'cancel']);
  const route = useRoute();

  const checkedList = reactive<Record<number, boolean>>(
    props.fileList.reduce((obj, file) => {
      obj[file.id] = false;
      return obj;
    }, {})
  );
  const renameIpt = ref<HTMLInputElement | null>(null);

  // 文件双击事件
  const fileClick = (index: number) => emit('fileClick', props.fileList[index]);
  // 文件点击
  const fileChoosed =  (index: number) => {
    checkedList[index] = !checkedList[index];
  }
  // 右键点击时勾选，事件传播到父组件，弹出右键菜单
  const handleContextMenu = (e: MouseEvent, id: number) => {
    e.preventDefault();
    checkedList[id] = true;
  };
  // 确认重命名
  const finishRenameOrCreate = async (file: formatFile) => {
    // 为了不打破单向数据流
    emit('finish', file, renameIpt.value[0].value);
    checkedList[file.id] = false;
  };
  // 取消重命名
  const cancelRename = (file: formatFile) => {
    emit('cancel', file);
    checkedList[file.id] = false;
  };
  // 勾选的文件列表变化时触发事件
  watch(checkedList, () => emit('fileChoosed', checkedList));
  // 文件刷新时更新checkedList
  watch(
    () => props.fileList,
    (newFileList) => {
      newFileList.forEach((file) => {
        if (checkedList[file.id] === undefined) {
          checkedList[file.id] = false;
        }
      });
    }
  );
  // 切换目录时清空 checkedList
  watch(
    () => route.query.path,
    () => {
      Object.keys(checkedList).forEach((key) => {
        delete checkedList[key];
      });
    }
  );

  const clearSelection = () => Object.keys(checkedList).map((id) => (checkedList[id] = false));
  defineExpose({ clearSelection });
</script>

<style lang="scss">
  .file-item {
    width: 10rem;
    padding: 10px;
    margin-bottom: 1.5rem;
    text-align: center;
    vertical-align: top;
    border-radius: 15px;
    position: relative;
    &:hover,
    &.checked {
      background-color: var(--ep-file-hover);
      transform: scale(1.01);
      transition: transform 0.2s ease-out;
      .file-checkbox {
        visibility: visible;
      }
    }
    .file-checkbox {
      display: block;
      visibility: hidden;
      margin-bottom: 10px;
      width: 0.875rem;
      height: 0.875rem;
    }
    .picture {
      width: 8rem;
      height: 8rem;
      border-radius: 10px;
    }
    .rename-ipt {
      width: 100%;
      height: 23px;
      line-height: 23px;
      font-size: 13px;
      padding: 0 10px;
      border-radius: 4px;
      border: 1px solid #a7adbc;
      outline: none;
    }
    .rename-btn {
      width: 0.875rem;
      height: 0.875rem;
      padding: 2px;
      top: 0.7rem;
      right: 0.7rem;
      position: absolute;

      &:first-of-type {
        right: 2rem;
      }
    }
    .filename {
      width: 100%;
      white-space: nowrap;
      text-overflow: ellipsis;
      font-size: 14px;
      line-height: 18px;
      color: var(--ep-table-text-color);
      overflow: hidden;
    }
    .describe {
      font-size: 12px;
      color: #818999;
      line-height: 16px;
      margin: 8px 0;
    }
  }
</style>
