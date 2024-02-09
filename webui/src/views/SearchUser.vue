<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			sessionUsername: sessionStorage.getItem("username"),
			sessionUserID: sessionStorage.getItem("user-id"),
			sessionToken: sessionStorage.getItem("auth-token"),
            username: "",
            userSums: [],
            search: false,
		}
	},
	methods: {
        async refresh(){
            this.loading = true;
            this.errormsg = false;
            this.username = "";
            this.search = false;
            this.userSums = [];
            this.loading = false;
        },

		async listUsers(){
            this.loading = true;
            this.errormsg = false;
            this.search = true;
            this.userSums = [];
            try {
                let response = await this.$axios.get("/users/", {
                    params: {
                        search: this.username, 
                    },

                    headers: {
                        Authorization: this.sessionToken,
                    }
                })

                if (response != null){
                    for(let i = 0; i < response.data.length; i++){
                        let userSum = {
                            userID: response.data[i]['user-id'],
                            username: response.data[i]['username'],
                            image: this.getImageFile(response.data[i]['profile-image-path'])
                        };
                        this.userSums.push(userSum);
                    }
                }

            } catch (e){
                if (e.response != null){
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        async getImageFile(imagePath) {
                this.loading = true;
                this.errormsg = null;
                let image = null;
                try {
                    let response = await this.$axios.get("/images/" + imagePath, {
                            headers: {
                                Authorization: this.sessionToken,
                            }
                        }
                    );
                    let ext = response.headers['content-type'].split('/')[1];
                    image = 'data:image/'+ext+';base64,'+response.data;
                } catch(e) {
                    if (!e.response.status == 404){
                        this.errormsg = e.toString();
                    }
                }
                this.loading = false;
                return image
            },
	},
	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Search</h1>
        <div class="btn-toolbar mb-2 mb-md-0">
            <div class="btn-group me-2">
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
                    Refresh
                </button>
            </div>
        </div>
    </div>
    <div id="search-box">
        <div>
            <form @submit.prevent="listUsers">
                <label for="username" class="form-label">Find a user using the username</label>
                <br>
                <input
                    id="username"
                    v-model="username"
                    type="text"
                    placeholder="e.g. Maria"
                />
                <button type="submit" title="Search"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg></button>
            </form>
        </div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
    <hr>
    <div v-if="search">
        <h4>Users</h4>
        <div id="user-list-box" v-if="userSums.length > 0">
            <ul>
                <div class="user-sum-box" v-for="(userSum, index) in userSums" :key="index">
                    <li>
                        <RouterLink :to="'/users/' + userSum.username + '/'" class="nav-link" v-if="this.sessionUserID != userSum.userID">
                            <!--
                            <img :src="userSum.image" alt="Profile image" class="rounded-circle mb-3" style="width: 30px;" v-if="image != null">
                            <img v-else src="https://yourteachingmentor.com/wp-content/uploads/2020/12/istockphoto-1223671392-612x612-1.jpg" class="rounded-circle mb-3" style="width: 30px;">
                            -->
                            {{ userSum.username }}
                        </RouterLink>
                        <RouterLink v-else to='/my-profile/' class="nav-link">
                            <!--
                            <img :src="userSum.image" alt="Profile image" class="rounded-circle mb-3" style="width: 30px;" v-if="image != null">
                            <img v-else src="https://yourteachingmentor.com/wp-content/uploads/2020/12/istockphoto-1223671392-612x612-1.jpg" class="rounded-circle mb-3" style="width: 30px;">
                            -->
                            {{ userSum.username }}
                        </RouterLink>
                    </li>
                </div>
                
            </ul>
            
        </div>
        <div v-else>
            <!--
                <h5><img src="https://i.redd.it/4s978dxj7xp51.jpg" style="width: 100px; heigth: 100px;"></h5>
            -->
            <p>...</p>
        </div>
	</div>
    <div v-else>
        <h5>Here you'll see the list of users</h5>
    </div>
</template>

<style>
</style>
