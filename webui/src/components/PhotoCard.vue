<script>
    export default {
        props: ['photo'],
        data: function() {
            return {
                errormsg: null,
                loading: false,
                username: sessionStorage.getItem("username"),
                authToken: sessionStorage.getItem("auth-token"),
                userID: sessionStorage.getItem("user-id"),
                image: "",
            }
        },

        methods: {
            async getImageFile() {
                this.loading = true;
                this.errormsg = null;
                try {
                    let response = await this.$axios.get("/images/" + this.photo.imagePath, {
                            headers: {
                                Authorization: this.authToken,
                            }
                        }
                    );
                    let ext = response.headers['content-type'].split('/')[1];
                    this.image = 'data:image/'+ext+';base64,'+response.data;
                } catch(e) {
                    this.errormsg = e.toString;
                }
                this.loading = false;
            },

            async likePhoto() {
                this.loading = true;
                this.errormsg = null;
                try {
                    await this.$axios.put("/photos/" + this.photo.photoID + "/likes/" + this.userID, {},{
                            headers: {
                                Authorization: this.authToken,
                            }
                        }
                    );

                } catch (e) {
                    this.errormsg = e.toString();
                }
                this.loading = false;
            },
            
            async unlikePhoto() {
                this.loading = true;
                this.errormsg = null;
                try {
                    await this.$axios.delete("/photos/"+this.photo.photoID+"/likes/"+this.userID, {},{
                        headers: {
                            Authorization: this.authToken,
                        }
                    });

                } catch (e) {
                    this.errormsg = e.toString();
                }
                this.loading = false;
            },
            
            async deletePhoto(){
                this.loading = true;
                this.errormsg = false;
                try {
                    let response = await this.$axios.delete("/photos/" + this.photo.photoID + "/", {
                        headers : {
                            Authorization: this.authToken
                        }
                    })
                    console.log(response.status)
                } catch (e) {
                    this.errormsg = e.toString();
                }

                this.loading = false;
            },
        },

        mounted() {
            this.getImageFile();
        }
    }
</script>

<template>
        <div class="row">
            <div class="col-md-3">
                <div class="card mb-3 shadow-sm">
                    <div class="card-header">
                        <div class="d-flex justify-content-between align-items-center mt-2">
                            <h2 class="text-muted">{{photo.username}}</h2>
                            <button @click="deletePhoto" v-if="photo.ownerID == userID" class="btn btn-outline-primary btn-sm">&times;</button>
                        </div>
                    </div>
                    <img :src="image" alt="Photo" class="card-img-top">
                    <div class="card-body">
                        <div class="d-flex align-items-right">
                            <small class="text-muted">{{ photo.timestamp }}</small>
                        </div>
                        <div>
                            <p>{{ photo.caption }}</p>
                        </div>
                        <div class="d-flex justify-content-between align-items-center mt-2">
                            <span class="text-muted">{{ photo.likesCount }} likes</span>
                            <span class="text-muted">{{ photo.commentsCount }} comments</span>
                        </div>
                        <div class="d-flex justify-content-between align-items-center mt-2">
                            <button @click="likePhoto" class="btn btn-outline-primary btn-sm">Like</button>
                            <button @click="unlikePhoto" class="btn btn-outline-primary btn-sm">Unlike</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

</template>

<style>
</style>