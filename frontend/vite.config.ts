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
      https: false, // 不启用HTTPS
      cors: true, // 启用CORS
      proxy: { // 代理配置
        "/devApi": {
          target: "http://104.156.140.42:18888", // 目标服务器
          changeOrigin: true, // 修改请求头中的origin为目标URL
          rewrite: path => path.replace(/^\/devApi/, "") // 路径重写
        },
      }
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
      assetsInlineLimit: 8 * 1024, // 小于8KB的资源转为base64
      sourcemap: isProd, // 生产环境生成sourcemap
      outDir: env.VITE_OUTDIR, // 输出目录（从环境变量获取）
      emptyOutDir: true, // 构建前清空输出目录
      chunkSizeWarningLimit: 1500, // 块大小警告阈值(KB)
      rollupOptions: { // Rollup打包配置
        output: {
          manualChunks(id) { // 手动分块策略
            if (id.includes('node_modules')) {
              // 将node_modules中的依赖单独分块
              return id.toString().split('node_modules/')[1].split('/')[0].toString();
            }
          },
          // 输出文件命名规则（添加时间戳避免缓存）
          chunkFileNames: `static/js/[name]-[hash]${Timestamp}.js`,
          entryFileNames: `static/js/[name]-[hash]${Timestamp}.js`,
          assetFileNames: `static/[ext]/[name]-[hash]${Timestamp}.[ext]`,
        }
      }
    },

    // JSON文件处理配置
    json: {
      stringify: false
    }
  }
})