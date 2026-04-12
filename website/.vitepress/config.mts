import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'pipectl',
  description: 'Run YAML-defined data pipelines from the command line.',
  base: '/pipectl/',
  srcDir: './docs',

  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
  ],

  themeConfig: {
    logo: {
      light: '/logo-light.svg',
      dark: '/logo-dark.svg',
      alt: 'pipectl',
    },
    siteTitle: false,

    nav: [
      { text: 'Guide', link: '/getting-started' },
      { text: 'Steps', link: '/steps/' },
      { text: 'Examples', link: '/examples/' },
      {
        text: 'GitHub',
        link: 'https://github.com/pipectl/pipectl',
      },
    ],

    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'Core Concepts', link: '/concepts' },
          { text: 'CLI Reference', link: '/cli' },
          { text: 'Payload Formats', link: '/formats' },
        ],
      },
      {
        text: 'Steps',
        items: [
          { text: 'Overview', link: '/steps/' },
          { text: 'normalize', link: '/steps/normalize' },
          { text: 'filter', link: '/steps/filter' },
          { text: 'redact', link: '/steps/redact' },
          { text: 'cast', link: '/steps/cast' },
          { text: 'default', link: '/steps/default' },
          { text: 'select', link: '/steps/select' },
          { text: 'rename', link: '/steps/rename' },
          { text: 'sort', link: '/steps/sort' },
          { text: 'limit', link: '/steps/limit' },
          { text: 'dedupe', link: '/steps/dedupe' },
          { text: 'convert', link: '/steps/convert' },
          { text: 'validate-json', link: '/steps/validate-json' },
          { text: 'assert', link: '/steps/assert' },
          { text: 'http-transform', link: '/steps/http-transform' },
          { text: 'log', link: '/steps/log' },
          { text: 'count', link: '/steps/count' },
        ],
      },
      {
        text: 'Examples',
        items: [
          { text: 'Gallery', link: '/examples/' },
          { text: 'Customer Signup', link: '/examples/customer-signup' },
          { text: 'CSV Intake', link: '/examples/csv-intake' },
          { text: 'Audit Export', link: '/examples/audit-export' },
          { text: 'Service-to-Service', link: '/examples/service-to-service' },
        ],
      },
      {
        text: 'Contributing',
        items: [{ text: 'Contributing Guide', link: '/contributing' }],
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/pipectl/pipectl' },
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © pipectl contributors',
    },

    search: {
      provider: 'local',
    },
  },
})
