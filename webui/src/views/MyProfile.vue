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
			modalProPic: false,
            newUsername: "",
			uploadPhotoFile: null,
			image: null,
		}
	},
	methods: {

		refresh(){
			this.profile = {};
			this.photos = [];
			this.newUsername = "";
			this.errormsg = null;
			this.errormsgChUname = null;
			this.getUserProfile(this.username);
			this.getUserPhotos(this.username);
		},

		async getImageFile(imagePath){
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
				this.image = 'data:image/'+ext+';base64,'+response.data;
			} catch(e) {
				if (e.response.status != 404){
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},

		async getUserProfile() {
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
				this.getImageFile(this.profile['profile-image-path']);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;	
		},

		async getUserPhotos(){
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.username + "/photos/", {
						headers: {
							Authorization: this.authToken,
						}
					}
				);
				if (response.data != null){
					for(let i=0; i<response.data.length; i++){
						let photoID = response.data[i]['photo-id'];
						let ownerID = response.data[i]['owner-id'];
						let ownerUsername = response.data[i]['owner-username'];
						let timestamp = response.data[i]['timestamp'];
						let imagePath = response.data[i]['image-path'];
						let likesCount = response.data[i]['likes-count'];
						let commentsCount = response.data[i]['comments-count'];
						let caption = response.data[i]['caption'];
						let photo = {photoID, ownerID, ownerUsername, timestamp, imagePath, likesCount, commentsCount, caption};
						this.photos.push(photo);
					}
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
				this.refresh();
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

		showModalProPic(){
			this.modalProPic = true;
		},

		hideModalProPic(){
			this.modalProPic = false;
		},

		uploadFile(){
			this.uploadPhotoFile = this.$refs.file.files[0];
		}, 

		async uploadProfileImage(){
			this.loading = true;
			this.errormsg = null;

			let ftype = this.$refs.file.files[0].name.split('.');
			ftype = ftype[ftype.length - 1];
			if (ftype == "jpeg"){
				ftype = "jpg";
			}

			if (this.uploadPhotoFile == null){
				this.errormsg = "You must select an image file (png/jpg)";
				this.uploadPhotoFile = null;
				this.loading = false;
				return null;
			}
			if (this.uploadPhotoFile.size/1024 > 16500){
				// Check if the file size is greater than 16 MB (16500 kB)
				this.errormsg = "Maximum file size: 16 MB";
				this.uploadPhotoFile = null;
				this.loading = false;
				return null;
			}
			const headers = { 'Content-Type': 'image/'+ftype, 
									'Authorization': this.authToken,
									'Access-Control-Allow-Origin': '*'};
			try{
				let response = await this.$axios.put("/users/"+this.username+"/profile-image", this.uploadPhotoFile, { headers });
				this.hideModalPropic();
				this.refresh();
			} catch (e){
				this.errormsg = e.toString();
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
		<h1 class="h2">My Profile</h1>
		<div class="btn-toolbar mb-2 mb-md-0">
			<div class="btn-group me-2">
				<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
					Refresh
				</button>
			</div>
		</div>
	</div>
	<div v-if="!loading">
		<div class="container-fluid">
			<div class="row">
				<div class="col d-flex align-items-center">
					<div v-if="image != null">
					<img :src="image" alt="Profile image" style="width: 130px; height: 125px; border-radius: 50%; object-fit: cover; border: 1px solid #000;"/>
					</div>
					<div v-else>
						<img src="https://yourteachingmentor.com/wp-content/uploads/2020/12/istockphoto-1223671392-612x612-1.jpg" style="width: 130px; height: 125px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
					</div>
					<div v-if="!modalProPic">
						<button @click="showModalProPic" class="btn btn-icon btn-sm"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#tool"/></svg></button>
					</div>
					<div class="upload-propic-box" v-else>
						<button @click="hideModalProPic" class="btn btn-icon btn-sm">&times;</button>
						<form @submit.prevent="uploadProfileImage">
							<label for="uploadProfileImage">Select a profile image</label>
							<input
								id="uploadProfileImage"
								@change="uploadFile"
								type="file"
								ref="file"
								accept = ".png, .jpg, .jpeg"
							>
							<button type="submit">Upload</button>
						</form>
					</div>
				</div>
				<div id="user-stats" class="col">
					<ul class="list-group list-group-horizontal">
						<li class="list-group-item text-center w-50">
							Photos
							<br>
							{{ profile['photos-count'] }}
						</li>
						<li class="list-group-item text-center w-50">
							Followers
							<br>
							{{ profile['followers-count'] }}
						</li>
						<li class="list-group-item text-center w-50">
							Following
							<br>
							{{ profile['followings-count'] }}
						</li>
					</ul>
				</div>
				<div class="col"></div>
			</div>
		</div>
		<br>
		<div class="container-fluid">
			<div class="row">
				<div id="username-box" class="col d-flex align-items-center">
					<h5>Username: {{ username }}</h5>
					<div v-if="!buttonModal">
						<button class="btn btn-icon btn-sm" @click="showButtonModal" title="Change username"><svg class="feather"><use href="/feather-sprite-v4.29.0.svg#tool"/></svg></button>
					</div>
					<div v-else>
						<h6>
							Change username
							<button class="btn btn-icon btn-sm" @click="hideButtonModal">&times;</button>
						</h6>
						<div>
							<form @submit.prevent="setMyUserName">
								<label for="newUsername">New Username</label>
								<br>
								<input
									id="newUsername"
									v-model="newUsername"
								/>
								<button class="btn btn-sm btn-primary" type="submit">Confirm</button>
							</form>
						</div>
					</div>
					<ErrorMsg v-if="errormsgChUname" :msg="errormsgChUname"></ErrorMsg>
				</div>
			</div>
		</div>

    	<hr>
		<div class="profile-photos">
			<h3>Photos</h3>
    		<PhotoCard v-for="photo in photos" :key="photo.photoID" :photo="photo" v-if="photos.length!=0" @photoUpdated="refresh"/>
			<div v-else>
				<h5>There are no photos yet :'(</h5>
				<!--
					<h6>Post something<img src="https://i.redd.it/4s978dxj7xp51.jpg" style="width: 100px; heigth: 100px;"></h6>
				-->
			</div>
    	</div>
	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>
