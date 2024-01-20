<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: null,
			userID: null,
			token: null,
		}
	},
	methods: {
		async doLogin() {
			this.loading = true;
			this.errormsg = null;
			if (this.username == null) {
				this.errormsg = "The username can not be empty";
			} else if (this.username.length < 3 || this.username.length > 16) {
				this.errormsg = "The username must have a minimum of 3 characters and a maximum of 16 characters";
			} else {
				try {
					let response = await this.$axios.post("/login", { username: this.username });
					this.userID = response.data["user-id"];
					this.token = response.data["auth-token"];
					sessionStorage.setItem("user-id", this.userID);
					sessionStorage.setItem("auth-token", this.token);
					this.$router.push({ path: "/" });
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.loading = false;
		}
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<h1>Login page</h1>
		<div id="login-box">
			<span>
				<label for="username">Username</label>
				<br>
				<input type="text" id="username" v-model="username" placeholder="e.g. Maria" required>
				<button type="submit" @click="doLogin">Login</button>
			</span>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>

<style>
</style>
