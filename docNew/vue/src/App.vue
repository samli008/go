<template>
  <div id="app">
    <el-container style="height: 100vh; border: 1px solid #eee">
      <el-aside width="170px" style="background-color: rgb(238, 241, 246)">
          <el-menu router :default-openeds="open_list">
            <el-submenu index="1">
                <template slot="title">
                    <i class="el-icon-tickets"></i>Documents
                </template>
                <el-menu-item index="/linux">Linux</el-menu-item>
                <el-menu-item index="/dell">Dell</el-menu-item>
                <el-menu-item index="/netapp">Netapp</el-menu-item>
                <el-menu-item index="/vmware">Vmware</el-menu-item>
                <el-menu-item index="/windows">Windows</el-menu-item>
                <el-menu-item index="/private">Private</el-menu-item>
            </el-submenu>
          </el-menu>
      </el-aside>
      
      <el-container>
          <!-- <el-header style="text-align: right; font-size: 15px">
          <span v-if="username">欢迎:{{username}}</span><button style="margin-left: 10px; font-size: 15px;" @click="logout">注销</button>
          </el-header> -->
          
          <el-main>
              <!-- <keep-alive include="labPage">
                  <router-view :key="$route.fullPath"></router-view>
              </keep-alive> -->
              <router-view></router-view>
          </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'app',
  data(){
    return {
      open_list: ['1']
    }
  },
}
</script>

<style>
html,body{
    padding: 0;
    margin: 0;
}
.el-header {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
}

.el-aside {
    color: #333;
}
.login-box{
    width: 350px;
    margin: 80px auto;
    border: 1px solid #dcdfe6;
    padding: 40px;
    border-radius: 10px;
    box-shadow: 0 0 30px;
}
.login-title{
    text-align: center;
    margin-bottom: 20px;
}
.el-input{
    width: 92%;
}
.item-button{
    display: flex;
}
.item-button button{
    width: 120px;
    margin: 0 auto;
}
.userBox{
    width: 600px;
    height: 800px;
    margin: 100px auto;
    padding: 20px;
    /* border: 1px solid #ccc; */
    border-radius: 8px;
    box-shadow: 0 0 5px #eee;
}
</style>
