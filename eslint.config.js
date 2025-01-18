import antfu from '@antfu/eslint-config'

export default antfu({
  rules: {
    'antfu/no-top-level-await': 'off',
    'no-console': 'off',
    'style/max-statements-per-line': 'off',
    'node/prefer-global/process': 'off',
  },
})
