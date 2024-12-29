<template>
  <div class="play-container">
    <el-card class="play-card" shadow="always">
      <!-- close -->
      <div class="card-header">
        <h2>{{ playInfo.name }}</h2>
        <ElIcon size="16" class="close-btn" @click="hanleClose"><CloseBold /></ElIcon>
      </div>
      <!-- 内容 -->
      <div v-if="playInfo.type === 'video'" 
        ref="video" 
        class="inner-wrapper">
        <video-player
        :src=playInfo.url
        :controls="true"
        :loop="false"
        :volume="0.6"
        :autoplay="true"
        :speedRate="true"
        :playbackRates=[0.5,1,1.5,2,4,8]
        :fluid="true"
        language="zh-cn"
        preload="auto"
        class="inner-wrapper" />
      </div>
      <audio
        v-else-if="playInfo.type === 'audio'"
        :src="playInfo.url"
        controls
        class="inner-wrapper"></audio>
      <div v-else-if="playInfo.type === 'md'" class="md-content-container inner-wrapper">
        <div class="md-content" ref="mdHtml"></div>
      </div>
      <!-- 图片展示 -->
      <div v-else-if="playInfo.type === 'picture'" class="picture">
        <el-image
          :src="playInfo.url"
          :lazy="true"
          class="picture"
          fit="contain"
          preview-teleported
          @click.stop="() => {}" />
      </div>
      <!-- txt | pdf -->
      <iframe
        v-else
        class="inner-wrapper"
        title="preview"
        frameborder="0"
        height="600"
        :src="playInfo.url"></iframe>
    </el-card>
    <div class="panel-cover" v-show="playPageShow"></div>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, ref } from 'vue';
  import request from '@/utils/request';
  import { CloseBold } from '@element-plus/icons-vue';

  import MarkdownIt from 'markdown-it';
  import { full as emoji } from 'markdown-it-emoji';
  import Mark from 'markdown-it-mark';

  import { VideoPlayer } from '@videojs-player/vue'
  import 'video.js/dist/video-js.css'

  defineOptions({ name: 'PlayPage' });
  const props = defineProps<{
    playPageShow: boolean;
    playInfo: {
      name: string;
      url: string;
      type: string;
    };
  }>();
  const emit = defineEmits(['update:playPageShow']);

  const video = ref<HTMLElement>(null);
  const mdHtml = ref<HTMLElement>(null);

  const md = MarkdownIt({
    html: true, // 可以识别html
    linkify: true, // 自动检测像链接的文本
    typographer: true, // 优化排版，标点
    // 代码着色
    highlight: function (str, lang) {
      // @ts-ignore
      if (lang && hljs.getLanguage(lang)) {
        try {
          return (
            '<pre><code class="hljs">' +
            // @ts-ignore
            hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
            '</code></pre>'
          );
        } catch (__) {}
      }

      return '<pre><code class="hljs">' + md.utils.escapeHtml(str) + '</code></pre>';
    },
  });
  md.use(emoji).use(Mark);
  const initRender = async () => {
    if (props.playInfo.type === 'md') {
      const res = await request.get(props.playInfo.url+"?replaceImgPath=true");
      mdHtml.value.innerHTML = md.render(res);
    }
  };

  const hanleClose = () => emit('update:playPageShow', false);

  onMounted(() => {
    if (mdHtml.value) {
      initRender();
    }
  });

</script>

<style scoped lang="scss">
  .play-card {
    width: 60rem;
    border-radius: 10px;
    text-align: initial;
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 100;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.8rem;
      .close-btn {
        color: var(--ep-contextmenu-text-color);
        cursor: pointer;
        &:hover {
          color: var(--ep-menu-active-color);
        }
      }
    }
    .picture {
      width: 100%;
      height: calc(100vh - 9rem);
      text-align: center;
      border-radius: 10px;
    }
    .inner-wrapper {
      width: 100%;
      max-height: calc(100vh - 9rem);
      border-radius: 10px;
    }

    .md-content-container {
      padding: 20px 30px;
      font-size: 14px;
      outline: 1px solid rgb(218, 216, 216);
      overflow-y: auto;

      :deep() {
        h1 {
          padding-bottom: 0.4rem;
          font-size: 2.2rem;
          line-height: 1.3;
        }
        h2 {
          font-size: 1.75rem;
          line-height: 1.225;
          margin: 35px 0 15px;
          padding-bottom: 0.5em;
          border-bottom: 1px solid #ddd;
        }
        h3 {
          font-size: 1.4rem;
          line-height: 1.43;
          margin: 20px 0 7px;
        }
        h4 {
          font-size: 1.2rem;
        }
        h5 {
          font-size: 1rem;
        }
        h6 {
          font-size: 1rem;
          color: #777;
        }
        p {
          line-height: 1.6rem;
          word-spacing: 0.05rem;
        }
        pre .hljs {
          border-radius: 10px;
          background: #f9f5f5;
        }
        ol li {
          padding-left: 0.5rem;
        }
        ul,
        ol {
          padding-left: 30px;
        }
        p,
        blockquote,
        ul,
        ol,
        dl,
        table {
          margin: 0.8em 0;
        }

        li > ol,
        li > ul {
          margin: 0 0;
        }

        hr {
          height: 2px;
          padding: 0;
          margin: 16px 0;
          background-color: #e7e7e7;
          border: 0 none;
          overflow: hidden;
          box-sizing: content-box;
        }
        blockquote {
          border-left: 4px solid #42b983;
          padding: 10px 15px;
          color: #777;
          background-color: rgba(66, 185, 131, 0.1);
        }
        thead {
          background-color: #fafafa;
        }
        mark {
          background-color: #81a781;
          border-radius: 2px;
          padding: 2px 4px;
          margin: 0 2px;
          color: #222;
          font-weight: 500;
        }
        table {
          padding: 0;
          word-break: initial;
          tr {
            border-top: 1px solid #dfe2e5;
            margin: 0;
            padding: 0;
            &:nth-child(2n) {
              background-color: #fafafa;
            }
            th {
              font-weight: bold;
              border: 1px solid #dfe2e5;
              border-bottom: 0;
              text-align: left;
              margin: 0;
              padding: 6px 13px;
            }
            td {
              border: 1px solid #dfe2e5;
              text-align: left;
              margin: 0;
              padding: 6px 13px;
            }
          }
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
    z-index: 99;
    background-color: rgba(0, 0, 0, 0.5);
  }
</style>
