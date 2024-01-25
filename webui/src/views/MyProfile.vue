<script>
import PhotoCard from '../components/PhotoCard.vue';

export default {
	components: {
		PhotoCard
	},

	data: function() {
		return {
			errormsg: null,
			loading: false,
			username: sessionStorage.getItem("username"),
			userID: sessionStorage.getItem("user-id"),
			authToken: sessionStorage.getItem("auth-token"),
			profile: {},
			photos: [],
		}
	},
	methods: {
		async getMyUserProfile() {
			this.loading = true;
			this.errormsg = null;
			if (this.authToken == null) {
				this.$router.push({ path: "/login" });
			} else {
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
				try {
					let response = await this.$axios.get("/users/" + this.username + "/photos/", {
							headers: {
								Authorization: this.authToken,
							}
						}
					);
					for(let i=0; i<response.data.length; i++){
						let photoID = response.data[i]['photo-id'];
						let ownerID = response.data[i]['owner'];
						let timestamp = response.data[i]['timestamp'];
						let imagePath = response.data[i]['image-path'];
						let likesCount = response.data[i]['likes-count'];
						let commentsCount = response.data[i]['comments-count'];
						let caption = response.data[i]['caption'];
						let photo = {photoID, ownerID, timestamp, imagePath, likesCount, commentsCount, caption};
						this.photos.push(photo);
					}
				} catch (e) {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
	
	},
	mounted() {
		this.getMyUserProfile();
	}
}
</script>

<template>
	<div class="user-profile" v-if="!this.loading">
		<h1>{{ username }}'s profile page</h1>
		<div class="profile-image">
			<img :src="profile['profile-image-path']" alt="Profile image" />
			<button @click="editProfileImage">Edit profile image</button>
		</div>
		
		<div class="profile-stats">
    		<span>Photos: {{ profile['photos-count'] }}</span>
			<br>
    		<span>Following: {{ profile['followings-count'] }}</span>
			<br>
    		<span>Followers: {{ profile['followers-count'] }}</span>
    	</div>
    	
		<div class="profile-photos">
    		<PhotoCard v-for="photo in photos" :key="photo.photoID" :photo="photo"/>
    	</div>
	    
		<div class="profile-options">
    		<button @click="changeUsername">Change username</button>
    		<button @click="viewFollowers">View followers</button>
    		<button @click="viewFollowing">View following</button>
    		<button @click="viewBlocked">View blocked users</button>
    	</div>
	</div>
</template>

<style>
</style>
