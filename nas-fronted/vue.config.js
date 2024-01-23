module.exports = ({
  lintOnSave: false,

  devServer: {
    // open: true,
    proxy: {
      '/': {
        ws: false,
        target: 'http://localhost:8001', // 后端服务的地址
        changeOrigin: true,
      },
    },
  },
})
