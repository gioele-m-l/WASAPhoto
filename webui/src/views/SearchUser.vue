<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: sessionStorage.getItem("username"),
			userID: sessionStorage.getItem("user-id"),
			token: sessionStorage.getItem("auth-token"),
            username: "",
            userSums: [],
		}
	},
	methods: {
        async refresh(){
            this.loading = true;
            this.errormsg = false;
            this.username = "";
            this.userSums = [];
            this.loading = false;
        },
		async listUsers(){
            this.loading = true;
            this.errormsg = false;
            this.userSums = [];
            try {
                let response = await this.$axios.get("/users/", {
                    params: {
                        search: this.username, 
                    },

                    headers: {
                        Authorization: this.token,
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
                                Authorization: this.token,
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
                <label for="username" class="form-label">Username</label>
                <br>
                <input
                    id="username"
                    v-model="username"
                    type="text"
                    placeholder="e.g. Maria"
                />
                <button type="submit"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg></button>
            </form>
        </div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
    <hr>
    <div>
        <h4>Users</h4>
        <div id="user-list-box" v-if="userSums.length > 0">
            <div class="user-sum-box" v-for="userSum in userSums" :key="userSum.userID">
                <img :src="userSum.image" alt="Profile image" class="rounded-circle mb-3" style="width: 200px;">
                <RouterLink :to="'/users/' + userSum.username + '/'" class="nav-link" v-if="this.userID != userSum.userID">
					{{ userSum.username }}
				</RouterLink>
                <RouterLink v-else to='/my-profile/' class="nav-link">
					{{ userSum.username }}
				</RouterLink>
                <br>
            </div>
        </div>
        <div v-else>
            <h5><img src="https://i.redd.it/4s978dxj7xp51.jpg" style="width: 100px; heigth: 100px;"></h5>
        </div>
	</div>
</template>

<style>
</style>
