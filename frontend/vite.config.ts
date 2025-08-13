import { ConfigEnv, UserConfig, defineConfig, loadEnv } from 'vite'
import viteCompression from 'vite-plugin-compression'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { visualizer } from 'rollup-plugin-visualizer' // 打包分析插件

// 获取当前时间戳，用于构建输出的文件名
const Timestamp = new Date().getTime();

// 路径解析辅助函数
function _resolve(dir: string) {
  return path.resolve(__dirname, dir);
}

export default defineConfig(({ mode }: ConfigEnv): UserConfig => {
  // 判断当前是否为生产环境
  const isProd = mode === 'production'
  // 加载环境变量文件(.env)
  const env = loadEnv(mode, process.cwd())

  const enableVisualizer = isProd && env.VITE_VISUALIZER === 'true'

  return {
    // 基础公共路径，从环境变量VITE_ADDRESS_BASE_URL获取
    base: env.VITE_ADDRESS_BASE_URL,

    // CSS预处理器配置
    css: {
      preprocessorOptions: {
        scss: {
          api: 'modern-compiler', // 使用现代Sass编译器API
          silenceDeprecations: ['legacy-js-api'],
        }
      }
    },

    // 插件配置
    plugins: [
      // Vue插件
      vue(),

      // Gzip压缩插件（仅生产环境启用）
      viteCompression({
        verbose: true, // 显示压缩日志
        disable: !isProd, // 非生产环境禁用
        deleteOriginFile: false, // 不删除源文件
        threshold: 10240, // 文件大小10KB才压缩
        algorithm: 'gzip', // 压缩算法
        ext: '.gz', // 压缩文件扩展名
      }),

      // 打包分析插件（仅在显式启用时生效；容器/CI 默认不打开浏览器）
      enableVisualizer && visualizer({ open: false })
    ].filter(Boolean), // 过滤掉false的插件

    // 开发服务器配置
    server: {
      host: '0.0.0.0', // 监听所有IP
      port: 8080, // 端口号
      open: true, // 自动打开浏览器
      // https: false, // 不启用HTTPS
      cors: true, // 启用CORS
      // 开发环境热更新优化
      hmr: {
        overlay: true
      },
      // 代理配置 (已直接配置 API 基础地址，不需要代理)
      // proxy: { 
      //   "/": {
      //     target: "http://104.156.140.42:18888/", 
      //     changeOrigin: true, 
      //   },
      // }
    },

    // 依赖优化配置
    optimizeDeps: {
      include: [
        'vue',
        'vue-router',
        'pinia',
        'element-plus',
        '@element-plus/icons-vue',
        'axios',
        'vue-i18n'
      ],
      exclude: ['@vueuse/core'] // 排除某些包的预构建
    },

    // 模块解析配置
    resolve: {
      alias: { // 路径别名
        '@': _resolve('src'), // 源码目录
        '@assets': _resolve('src/assets'), // 资源目录
        '@views': _resolve('src/views'), // 视图目录
        '@components': _resolve('src/components'), // 组件目录
        '@utils': _resolve('src/utils'), // 工具函数
        '@router': _resolve('src/router'), // 路由配置
        '@store': _resolve('src/store'), // 状态管理
      }
    },

    // ESBuild配置
    esbuild: {
      pure: isProd ? ["console.log", "debugger"] : [] // 生产环境移除console.log和debugger
    },

    // 构建配置
    build: {
      minify: 'terser', // 使用terser进行代码压缩
      assetsInlineLimit: 4 * 1024, // 小于4KB的资源转为base64（优化）
      sourcemap: !isProd, // 开发环境生成sourcemap，生产环境不生成
      outDir: env.VITE_OUTDIR || 'dist', // 输出目录
      emptyOutDir: true, // 构建前清空输出目录
      chunkSizeWarningLimit: 1000, // 降低块大小警告阈值至1MB
      cssCodeSplit: true, // CSS代码分割
      terserOptions: {
        compress: {
          drop_console: isProd, // 生产环境移除console
          drop_debugger: isProd, // 生产环境移除debugger
          pure_funcs: isProd ? ['console.log'] : [] // 移除纯函数调用
        }
      },
      rollupOptions: { // Rollup打包配置
        output: {
          // 优化的手动分块策略
          manualChunks: {
            // Vue核心
            vue: ['vue', 'vue-router', 'pinia'],
            // Element Plus相关
            elementPlus: ['element-plus'],
            // Element Plus图标
            elementIcons: ['@element-plus/icons-vue'],
            // 工具库
            utils: ['axios'],
            // 其他第三方库单独打包
            vendor: ['vue-i18n']
          },
          // 输出文件命名规则（优化缓存策略）
          chunkFileNames: `static/js/[name].[hash].js`,
          entryFileNames: `static/js/[name].[hash].js`,
          assetFileNames: (assetInfo) => {
            const name = assetInfo.name || ''
            if (/\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/i.test(name)) {
              return `static/media/[name].[hash][extname]`
            } else if (/\.(png|jpe?g|gif|svg|ico|webp)(\?.*)?$/i.test(name)) {
              return `static/images/[name].[hash][extname]`
            } else if (/\.(woff2?|eot|ttf|otf)(\?.*)?$/i.test(name)) {
              return `static/fonts/[name].[hash][extname]`
            } else {
              return `static/assets/[name].[hash][extname]`
            }
          }
        },
        // 外部化依赖（CDN优化，可选）
        external: isProd ? [] : []
      }
    },

    // JSON文件处理配置
    json: {
      stringify: false
    }
  }
})