<script>
import UserCard from '../components/UserCard.vue';

export default {
    components: {
        UserCard,
    },

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
        logout(){
			sessionStorage.clear();
			this.$router.replace({ path: "/login"});
		},

        refresh(){
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
                        let user = {
                            userID: response.data[i]['user-id'],
                            username: response.data[i]['username'],
                            imagePath: response.data[i]['profile-image-path'],
                        };
                        this.userSums.push(user);
                    }
                }

            } catch (e){
                if (e.response != null){
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
	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Search</h1>
        <div class="btn-toolbar mb-2 mb-md-0">
            <div class="btn-group me-2">
                <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
                    Refresh
                </button>
            </div>
            <div class="btn-group me-2">
				<button type="button" class="btn btn-sm btn-outline-danger" @click="logout">
					Logout
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
        <div v-if="userSums.length > 0" class="w-25" id="user-list-box">
            <UserCard v-for="user in userSums" :key="user.userID" :user="user"></UserCard>
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
