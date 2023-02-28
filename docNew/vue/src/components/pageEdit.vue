<template>
    <div>
        <el-form @submit.native.prevent="saveArticle" ref="form" :model="article" label-width="80px">
            <el-form-item label="文档名称" prop="name" style="width:300px">
                <el-input v-model="article.name"></el-input>
            </el-form-item>
            <el-form-item label="文档类型" prop="type" style="width:300px">
                <el-input v-model="article.type"></el-input>
            </el-form-item>
            <el-form-item label="文档内容" prop="content" style="width:1000px">
                <el-input type="textarea" :autosize="{minRows:5,maxRows:100}" v-model="article.content"></el-input>
            </el-form-item>
            <div class="item-button">
                <el-button type="primary" native-type="submit">保存</el-button>
                <el-button @click="cancel">取消</el-button>
            </div>
        </el-form>
    </div>
</template>

<script>
export default {
    name: 'pageEdit',
    data(){
        return {
            article: {}
        }
    },
    methods: {
        saveArticle(){
            this.$axios.put(this.$route.params.type+'/doc/'+this.$route.params.id, this.article).then(()=>{
                this.$message({
                    type: 'success',
                    message: '保存成功!'
                })
                this.$router.push(this.$route.params.type)
            })
        },
        cancel(){
            this.$router.push(this.$route.params.type)
        }
    },
    created(){
        this.$axios.get(this.$route.params.type+'/content/'+this.$route.params.id).then(res=>{
            this.article=res.data.doc
        })
    }
}
</script>