import { defineStore } from 'pinia'
import { isExternal } from '@/utils'
import { basePermissions } from '@/settings'
import api from '@/api'

const routeComponents = import.meta.glob('@/views/**/*.vue')

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    menus: [],
    accessRoutes: [],
    asyncPermissions: [],
  }),
  getters: {
    permissions() {
      return basePermissions.concat(this.asyncPermissions)
    },
  },
  actions: {
    async initPermissions() {
      console.log('initPermissions')
      const { data } = (await api.getRolePermissions()) || []
      this.asyncPermissions = data
      this.menus = this.permissions
        .filter((item) => item.type === 'MENU')
        .map((item) => this.getMenuItem(item))
        .filter((item) => !!item)
        .sort((a, b) => a.sortOrder - b.sortOrder)
    },
    getMenuItem(item, parent) {
      const route = this.generateRoute(item, item.show ? null : parent?.key)
      if (item.enable && route.path && !isExternal(route.path)) this.accessRoutes.push(route)
      if (!item.show) return null
      const menuItem = {
        label: route.meta.title,
        key: route.name,
        path: route.path,
        icon: () => h('i', { class: `${route.meta.icon}?mask text-16` }),
        sortOrder: item.sortOrder ?? 0,
      }
      const children = item.children?.filter((item) => item.type === 'MENU') || []
      if (children.length) {
        menuItem.children = children
          .map((child) => this.getMenuItem(child, menuItem))
          .filter((item) => !!item)
          .sort((a, b) => a.sortOrder - b.sortOrder)
        if (!menuItem.children.length) delete menuItem.children
      }
      return menuItem
    },
    generateRoute(item, parentKey) {
      return {
        name: item.code,
        path: item.path,
        redirect: item.redirect,
        component: routeComponents[item.component] || undefined,
        meta: {
          icon: item.icon,
          title: item.name,
          layout: item.layout,
          keepAlive: !!item.keepAlive,
          parentKey,
          btns: item.children
            ?.filter((item) => item.type === 'BUTTON')
            .map((item) => ({ code: item.code, name: item.name })),
        },
      }
    },
    resetPermission() {
      this.$reset()
    },
  },
})
