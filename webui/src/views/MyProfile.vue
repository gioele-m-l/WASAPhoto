<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: sessionStorage.getItem("username"),
			userID: sessionStorage.getItem("user-id"),
			authToken: sessionStorage.getItem("auth-token"),
			profile: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			
			try {
				let response = await this.$axios.get("/users/" + this.username + "/", {
						headers: {
							Authorization: this.authToken,
						}
					} 
				);
				
				this.profile = response.data;

			} catch (e) {
				this.errormsg = e.toString();
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
		<h1>{{ username }}'s profile page</h1>
		<div>
			
		</div>
	</div>
</template>

<style>
</style>
