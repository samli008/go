<template>
    <div>
        <p v-html="article"></p>
    </div>
</template>
<script>
import MarkdownIt from 'markdown-it'
export default {
    name: 'labView',
    data(){
        return {
            article: {}
        }
    },
    methods: {
        markdown(post){
            const md = new MarkdownIt();
            const result = md.render(post);
            return result;
        }
    },
    created(){
        this.$axios.get('/content/'+this.$route.params.id).then(res=>{
            this.article=this.markdown(res.data.doc.content)
        })
    }
}
</script>