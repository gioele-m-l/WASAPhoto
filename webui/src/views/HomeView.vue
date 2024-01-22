<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			authToken: sessionStorage.getItem("auth-token"),
			userID: sessionStorage.getItem("user-id"),
			responseStream: null,
		}
	},
	methods: {
		async refresh() {
			if (this.authToken == null) {
				this.$router.push({ path: "/login" })
			}
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/photos/", {
						headers: {
							Authorization: this.authToken,
						}
					});
				this.responseStream = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			console.log(history)
		},
	},
	async refresh (){

	},

	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
	<div v-if="responseStream != null">
		<li v-if="responseStream.photos != null" v-for="p in responseStream.photos" :key="p.photoID">
			{{ p['photo-id'] }}, {{ p.timestamp }}, {{ p.owner }}, {{ p['image-path'] }}, {{ p['likes-count'] }}, {{ p['comments-count'] }}, {{ p.caption }}
		</li>
	</div>
</template>

<style>
</style>
