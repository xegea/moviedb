new Vue({
    el: "#app",
    data: {
        search : "",
        totalres : 0,
        contentList: [],
        lastPage: true,
        page: 1
    },
    methods: {
        async next() {
            try {
                if (this.search.length < 2) return;
                this.page++;
                const res = await axios.get("/search/?q=" + this.search + "&p=" + this.page + "&c=es-es");
                this.totalres = res.data.MovieList?.length;
                this.contentList = res.data.MovieList;
                this.lastPage = res.data.LastPage;
            } catch (err) {
                this.page--;
                console.log(err)
            }
        },
        async back() {
            try {
                this.page--;
                const res = await axios.get("/search/?q=" + this.search + "&p=" + this.page + "&c=es-es");
                this.totalres = res.data.MovieList?.length;
                this.contentList = res.data.MovieList;
                this.lastPage = res.data.LastPage;
            } catch (err) {
                this.page++;
                console.log(err);
            }
        }
    },
    watch : {
        search: _.debounce(async function() {
            try {
                if (this.search.length < 2) return;
                const res = await axios.get("/search/?q=" + this.search + "&p=1&c=es-es");
                if (res.data.MovieList == null) {
                    this.contentList = [];
                    return;
                }
                this.contentList = res.data.MovieList;
                this.totalres = res.data.MovieList?.length;
                this.lastPage = res.data.LastPage;
                //this.getResult = res.data[0].Title;
                this.page = 1;
            } catch (err) {
                console.log(err);
            }
        }, 700)
    }
});
