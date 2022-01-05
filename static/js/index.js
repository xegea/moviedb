new Vue({
    el: "#app",
    data: {
        search : "",
        totalres : 0,
        contentList: [],
        limit: 10,
        busy: false
    },
    methods: {
        async loadMore() {
            try {
                if (this.search.length < 2 || this.totalres < 10) return;
                console.log("Adding 10 more data results");
                this.busy = true;
                const res = await axios.get("https://frozen-peak-67206.herokuapp.com/content/?q=" + this.search + "&p=" + this.contentList.length / 10);
                this.totalres = res.data?.length;
                this.contentList = this.contentList.concat(res.data);
                this.busy = false;
            } catch (err) {
                console.log(err)
            }
        }
    },
    watch : {
        search: _.debounce(async function() {
            try {
                if (this.search.length < 2) return;
                const res = await axios.get("https://frozen-peak-67206.herokuapp.com/content/?q=" + this.search);
                this.contentList = res.data;
                this.totalres = res.data?.length;
                //this.getResult = res.data[0].Title;
            } catch (err) {
                console.log(err)
            }
        }, 700)
    }
});
