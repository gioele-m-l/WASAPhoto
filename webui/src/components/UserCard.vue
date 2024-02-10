<script>
    export default {
        props: ['user'],
        data: function() {
            return {
                errormsg: null,
                loading: false,
                sessionUsername: sessionStorage.getItem("username"),
                sessionAuthToken: sessionStorage.getItem("auth-token"),
                sessionUserID: sessionStorage.getItem("user-id"),
                image: null,
            }
        },

        methods: {

            refresh(){
                this.getImageFile();
            },

            async getImageFile() {
                this.loading = true;
                this.errormsg = null;
                try {
                    let response = await this.$axios.get("/images/" + this.user.imagePath, {
                            headers: {
                                Authorization: this.sessionAuthToken,
                            }
                        }
                    );
                    let ext = response.headers['content-type'].split('/')[1];
                    this.image = 'data:image/'+ext+';base64,'+response.data;
                } catch(e) {
                    if(e.response.status != 404){
                        this.errormsg = e.toString();
                    }
                }
                this.loading = false;
            },
        },

        mounted() {
            this.refresh();
        }
    }
</script>

<template>
    <div class="row">
        <div class="col">
            <div class="card mb-2 shadow-sm w-100">
                <RouterLink v-if="sessionUserID != user.userID" :to="'/users/' + user.username + '/'" class="d-flex justify-content-left flex-wrap flex-sm align-items-center nav-link">
                    <div class="row">
                        <div class="col">
                            <div v-if="image != null">
                                <img :src="image" alt="Profile image" style="width: 40px; height: 40px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
                            </div>
                            <div v-else>
                                <img src="https://yourteachingmentor.com/wp-content/uploads/2020/12/istockphoto-1223671392-612x612-1.jpg" style="width: 40px;height: 40px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
                            </div>
                        </div>
                        <div class="col">
                            {{ user.username }}
                        </div>
                    </div>
                </RouterLink>
                <RouterLink v-else to='/my-profile/' class="d-flex justify-content-left flex-wrap flex-sm align-items-center nav-link">
                    <div class="row">
                        <div class="col">
                            <div v-if="image != null">
                                <img :src="image" alt="Profile image" style="width: 40px; height: 40px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
                            </div>
                            <div v-else>
                                <img src="https://yourteachingmentor.com/wp-content/uploads/2020/12/istockphoto-1223671392-612x612-1.jpg" style="width: 40px;height: 40px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
                            </div>
                        </div>
                        <div class="col">
                            {{ user.username }}
                        </div>
                    </div>
                </RouterLink>
            </div>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>