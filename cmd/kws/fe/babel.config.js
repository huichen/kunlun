module.exports = {
  plugins: [
    [
      'import',
      {
        libraryName: 'element-plus',
        customStyleName: name => {
          return `element-plus/lib/theme-chalk/${name}.css`
        }
      }
    ]
  ]
}
