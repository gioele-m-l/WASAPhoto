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
					sessionStorage.setItem("username", this.username);
					this.$router.push({ path: "/" });
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
		async refresh() {
			this.loading = true;
			this.username = null;
			this.userID = null;
			this.token = null;
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
	<div class="container">
		<div class="row justify-content-center">
			<div class="col-md-6">
				<h1 class="text-center mt-5">Login page</h1>
				<div id="login-box">
					<form @submit.prevent="doLogin" class="p-4 rounded">
						<div class="form-group">
							<label for="text-box" class="form-label">Username</label>
							<input
								id="text-box"
								type="text"
								v-model="username"
								class="form-control"
								placeholder="e.g. Maria"
							/>
						</div>
						<div class="form-group">
							<button type="submit" class="btn btn-primary btn-block">Login</button>
						</div>
						
					</form>
					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
</style>
