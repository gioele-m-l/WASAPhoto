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
                try {
                    let response = await this.$axios.get("/images/" + imagePath, {
                            headers: {
                                Authorization: this.authToken,
                            }
                        }
                    );
                    let ext = response.headers['content-type'].split('/')[1];
                    return 'data:image/'+ext+';base64,'+response.data;
                } catch(e) {
                    if (!e.response.status == 404){
                        this.errormsg = e.toString();
                    }
                }
                this.loading = false;
                return null
            },
	},
	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div>
		<h1>Search</h1>
		<div id="search-box">
			<span>
				<form @submit.prevent="listUsers">
                    <label for="username">Username</label>
                    <input
                        id="username"
                        v-model="username"
                        type="text"
                    />
                    <button type="submit"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg></button>
                </form>
            </span>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
        <hr>
        <h2>Users</h2>
        <div id="user-list-box" v-if="userSums.length > 0">
            <div class="user-sum-box" v-for="userSum in userSums" :key="userSum.userID">
                <img :src="userSum.image" alt="Profile image">
                <RouterLink :to="'/users/' + userSum.username + '/'" class="nav-link" v-if="this.userID != userSum.userID">
					{{ userSum.username }}
				</RouterLink>
                <RouterLink v-else to='/my-profile/' class="nav-link">
					{{ userSum.username }}
				</RouterLink>
                <br>
            </div>
        </div>
        <h5 v-else>We didn't found any user :'(</h5>
	</div>
</template>

<style>
</style>
