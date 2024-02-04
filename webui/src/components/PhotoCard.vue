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
            }
        },

        mounted() {
            this.getImageFile();
        }
    }
</script>

<template>
    <div class="photo-card">
        <img :src="image" alt="Photo">
        <div class="photo-info">
            <h2>{{ photo.ownerID }}</h2>
            <p>{{ photo.timestamp }}</p>
            <p>{{ photo.likesCount }} likes</p>
            <p>{{ photo.commentsCount }} comments</p>
            <button @click="likePhoto">Like</button>
            <button @click="unlikePhoto">Unlike</button>
            <button @click="deletePhoto" v-if="photo.ownerID == userID">&times;</button>
        </div>
    </div>
</template>

<style scoped>
.photo-card {
  width: 300px;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 10px;
  margin: 10px;
  box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2);
  transition: 0.3s;
}

.photo-card:hover {
  box-shadow: 0 8px 16px 0 rgba(0, 0, 0, 0.2);
}

.photo-card img {
  width: 100%;
  height: auto;
}
</style>