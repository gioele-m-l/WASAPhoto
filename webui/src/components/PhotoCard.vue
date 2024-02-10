<script>
import UserCard from './UserCard.vue';

    export default {
    props: ['photo'],
    data: function () {
        return {
            errormsg: null,
            loading: false,
            username: sessionStorage.getItem("username"),
            authToken: sessionStorage.getItem("auth-token"),
            userID: sessionStorage.getItem("user-id"),
            image: "",
            liked: false,
            showComments: false,
            comments: null,
            newComment: "",
            user: null,
        };
    },
    methods: {
        refresh() {
            this.listComments();
            this.getPhotoLike();
        },
        async getImageFile() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.get("/images/" + this.photo.imagePath, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                let ext = response.headers['content-type'].split('/')[1];
                this.image = 'data:image/' + ext + ';base64,' + response.data;
            }
            catch (e) {
                this.errormsg = e.toString;
            }
            this.loading = false;
        },
        async likePhoto() {
            this.loading = true;
            this.errormsg = null;
            try {
                await this.$axios.put("/photos/" + this.photo.photoID + "/likes/" + this.userID, null, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                this.$emit('photoUpdated');
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async unlikePhoto() {
            this.loading = true;
            this.errormsg = null;
            try {
                await this.$axios.delete("/photos/" + this.photo.photoID + "/likes/" + this.userID, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                this.$emit('photoUpdated');
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async deletePhoto() {
            this.loading = true;
            this.errormsg = false;
            try {
                let response = await this.$axios.delete("/photos/" + this.photo.photoID + "/", {
                    headers: {
                        Authorization: this.authToken
                    }
                });
                this.$emit('photoUpdated');
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async getPhotoLike() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.get("/photos/" + this.photo.photoID + "/likes/" + this.userID, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                if (response != null) {
                    this.liked = true;
                }
            }
            catch (e) {
                if (e.response.status == 404) {
                    this.liked = false;
                }
                else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        async listComments() {
            this.loading = true;
            this.errormsg = null;
            this.comments = null;
            try {
                let response = await this.$axios.get("/photos/" + this.photo.photoID + "/comments/", {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                this.comments = response.data;
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async commentPhoto() {
            this.loading = true;
            this.errormsg = null;
            if (this.newComment.length <= 0) {
                this.errormsg = "The comment can not be empty";
                this.newComment = "";
                this.loading = false;
                return;
            }
            try {
                let response = await this.$axios.post("/photos/" + this.photo.photoID + "/comments/", { "text": this.newComment }, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                this.newComment = "";
                this.$emit('photoUpdated');
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async uncommentPhoto(commentID) {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.delete("/photos/" + this.photo.photoID + "/comments/" + commentID, {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                this.$emit('photoUpdated');
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        
        async getUserProfile(){
            this.loading = true;
            this.errormsg = null;

            try {
                let response = await this.$axios.get("/users/"+this.photo.ownerUsername+"/", {
                    headers: {
                        Authorization: this.authToken,
                    }
                })

                if(response != null){
                    this.user = {
                        userID: this.photo.ownerID,
                        username: this.photo.ownerUsername,
                        imagePath: response.data['profile-image-path']
                    }
                        
                }

            } catch(e) {
                this.errormsg = e.toString();
            }
            this.loading = false;

        },
    },
    mounted() {
        this.getUserProfile();
        this.getImageFile();
        this.refresh();
    },
    components: { UserCard }
}
</script>

<template>
        <div class="row">
            <div class="col-md-3">
                <div class="card mb-3 shadow-sm">
                    <div class="card-header">
                        <div class="d-flex justify-content-between mt-2 align-items-center" v-if="user != null">
                            <UserCard :user="user"></UserCard>
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
                            <div v-if="!liked">
                                <button @click="likePhoto" class="btn btn-outline-primary btn-sm">Like</button>
                            </div>
                            <div v-else>
                                <button @click="unlikePhoto" class="btn btn-outline-primary btn-sm">Unlike</button>
                            </div>
                            <div>
                                <button @click="showComments = !showComments" class="btn btn-outline-primary btn-sm">Comments</button>
                            </div>
                        </div>
                        <div v-if="showComments">
                            <hr>
                            <div>
                                <h6>Comments</h6>
                                <div v-for="comment in comments" :key="comment['comment-id']">
                                    <p class="d-flex justify-content-left flex-wrap flex-sm align-items-center">
                                        <RouterLink :to="'/users/'+comment['owner-username']+'/'" v-if="comment['owner-id'] != userID" class="nav-link">{{ comment['owner-username'] }}</RouterLink>
                                        <RouterLink to="/my-profile/" v-else class="nav-link">{{ comment['owner-username'] }}</RouterLink>
                                        : {{ comment.text }}
                                        <button @click="uncommentPhoto(comment['comment-id'])" v-if="comment['owner-id'] == userID" title="Delete comment" class="btn btn-icon btn-sm">&times;</button>
                                    </p>
                                    
                                </div>
                            </div>
                            <div>
                                <form @submit.prevent="commentPhoto" class="p-4 rounded">
                                    <div class="form-group">
                                        <label for="text-box" class="form-label">Post a comment</label>
                                        <input
                                            id="text-box"
                                            type="text"
                                            v-model="newComment"
                                            class="form-control"
                                            placeholder="Add a comment"
                                        />
                                    </div>
                                    <div class="form-group">
                                        <button type="submit" class="btn btn-primary btn-sm">Publish</button>
                                    </div> 
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>