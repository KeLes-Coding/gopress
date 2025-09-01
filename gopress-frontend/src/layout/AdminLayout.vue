<script setup>
// 导入 Element Plus 图标
import {
  House,
  Document,
  Collection,
} from '@element-plus/icons-vue'
</script>

<template>
  <div class="admin-layout">
    <el-container style="height: 100vh;">
      <el-aside width="200px">
        <div class="sidebar-logo">
          GoPress
        </div>
        <el-menu
          :default-active="$route.path"
          class="el-menu-vertical-demo"
          router 
        >
          <el-menu-item index="/admin/dashboard">
            <el-icon><house /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/admin/posts">
            <el-icon><document /></el-icon>
            <span>文章管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/categories">
            <el-icon><collection /></el-icon>
            <span>分类管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header>
          <div class="header-content">
            <span>后台管理系统</span>
            <div>
              <el-avatar size="small"> Admin </el-avatar>
              <span style="margin-left: 8px;">admin</span>
            </div>
          </div>
        </el-header>

        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style scoped>
.admin-layout .el-aside {
  background-color: #545c64;
  color: #fff;
}
.admin-layout .el-header {
  background-color: #fff;
  color: #333;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.sidebar-logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
}
.el-menu {
  border-right: none;
  background-color: #545c64;
}
.el-menu-item {
  color: #fff;
}
.el-menu-item:hover {
  background-color: #434a50;
}
.el-menu-item.is-active {
  background-color: #409EFF !important; /* Element Plus 的 active 样式权重较高，需要用 !important 覆盖 */
  color: #fff !important;
}
.el-main {
  background-color: #f0f2f5;
}
</style>