<script>
import PhotoCard from '../components/PhotoCard.vue';

export default {
	components: {
    PhotoCard,
},

	data: function() {
		return {
			errormsg: null,
			errormsgChUname: null,
			loading: false,
			username: sessionStorage.getItem("username"),
			userID: sessionStorage.getItem("user-id"),
			authToken: sessionStorage.getItem("auth-token"),
			profile: {},
			photos: [],
			buttonModal: false,
            newUsername: "",
		}
	},
	methods: {

		async getUserProfile(username) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + username + "/", {
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

		async getUserPhotos(username){
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + username + "/photos/", {
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
			this.loading = false;
		},
		
		
		showButtonModal(){
			this.buttonModal = true;
		},

		hideButtonModal(){
			this.buttonModal = false;
		},

		async setMyUserName(){
			this.loading = true;
			this.errormsg = null;
			this.errormsgChUname = null;
			try {
				let response = await this.$axios.put("/users/"+ this.username + "/username", {"username": this.newUsername },
					{ headers : {
						Authorization: this.authToken,
					}
				})
				this.hideButtonModal();
				
				this.username = this.newUsername;
				sessionStorage.setItem("username", this.username);
			} catch (e) {
				if(e.response.status == 400){
					this.errormsgChUname = "This username is invalid"
				}
				else if (e.response.status == 403){
					this.errormsgChUname = "This username is already in use, try with another one"
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},
	},

	mounted() {
		this.getUserProfile(this.username);
		this.getUserPhotos(this.username);
	}
}
</script>

<template>
	<div class="user-profile" v-if="!loading">
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

		<div>
			<h4>{{ username }}</h4>
			<button @click="showButtonModal">Change username</button>
			<div v-if="buttonModal">
					<h6>Change username</h6>
					<button @click="hideButtonModal">&times;</button>
				<form @submit.prevent="setMyUserName">
					<label for="newUsername">New Username</label>
					<input
						id="newUsername"
						v-model="newUsername"
					/>
					<button type="submit">Confirm</button>
				</form>
			</div>
			<ErrorMsg v-if="errormsgChUname" :msg="errormsgChUname"></ErrorMsg>
		</div>
    	<hr>
		<div class="profile-photos">
			<h3>Photos</h3>
    		<PhotoCard v-for="photo in photos" :key="photo.photoID" :photo="photo" v-if="photos.length!=0"/>
			<h5 v-else>There are no photos yet :'(</h5>
    	</div>
	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>
