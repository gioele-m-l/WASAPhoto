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
            sessionUsername: sessionStorage.getItem("username"),
			sessionUserID: sessionStorage.getItem("user-id"),
			sessionAuthToken: sessionStorage.getItem("auth-token"),

            username: this.$route.params.username,
			profile: {},
			photos: [],

			followed: false,
			banned: false,
		}
	},
	methods: {

		async refresh(){
			this.banned = false;
			this.followed = false;
			this.profile = {};
			this.photos = [];
			this.checkBan();
			if(this.banned != true){
				this.checkFollow();
			}
			this.getUserProfile(this.username);
			this.getUserPhotos(this.username);
		},

		async checkBan(){
			this.loading = true;
			this.errormsg = false;

			try {
				let response = await this.$axios.get("/users/" + this.sessionUsername + "/banned/", {
					headers : {
						Authorization: this.sessionAuthToken,
					}
				})

				if(response.data != null){
					for(let i=0; i<response.data.length; i++){
						if(response.data[i]["username"] == this.username){
							this.banned = true;
							break;
						}
					}
					
				}
			} catch (e){
				this.errormsg = e.toString();
			}

			this.loading = false;
		},

		async checkFollow(){
			this.loading = true;
			this.errormsg = false;

			try {
				let response = await this.$axios.get("/users/" + this.sessionUsername + "/followings/", {
					headers : {
						Authorization: this.sessionAuthToken,
					}
				})

				if(response.data != null){
					for(let i=0; i<response.data.length; i++){
						if(response.data[i]["username"] == this.username){
							this.followed = true;
							break;
						}
					}
					
				}
			} catch (e){
				this.errormsg = e.toString();
			}

			this.loading = false;
		},

		async getUserProfile(username) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + username + "/", {
						headers: {
							Authorization: this.sessionAuthToken,
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
							Authorization: this.sessionAuthToken,
						}
					}
				);
				if (response.data != null){
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
				}
			} catch (e) {
                if(e.response.status != 404){
				    this.errormsg = e.toString();
                }
			}
			this.loading = false;
		},

        async followUser(){
            this.loading = true;
            this.errormsg = false;
            try {
                let response = await this.$axios.put("/users/"+this.sessionUsername+"/followings/"+this.profile['user-id'], null, {
                    headers : {
                        Authorization: this.sessionAuthToken,
                    }
                })
				this.refresh();
            } catch(e){
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

        async unfollowUser(){
            this.loading = true;
            this.errormsg = false;
            try {
                let response = await this.$axios.delete("/users/"+this.sessionUsername+"/followings/"+this.profile['user-id'], {
                    headers : {
                        Authorization: this.sessionAuthToken,
                    }
                })
				this.refresh();
            } catch(e){
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

        async banUser(){
            this.loading = true;
            this.errormsg = false;
            try {
                let response = await this.$axios.put("/users/"+this.sessionUsername+"/banned/"+this.profile['user-id'], null, {
                    headers : {
                        Authorization: this.sessionAuthToken,
                    }
                })
				this.refresh();
            } catch(e){
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

        async unbanUser(){
            this.loading = true;
            this.errormsg = false;
            try {
                let response = await this.$axios.delete("/users/"+this.sessionUsername+"/banned/"+this.profile['user-id'], {
                    headers : {
                        Authorization: this.sessionAuthToken,
                    }
                })
				this.refresh();
            } catch(e){
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
		<h1 class="h2">{{ profile['username'] }}'s profile</h1>
		<div class="btn-toolbar mb-2 mb-md-0">
			<div class="btn-group me-2">
				<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
					Refresh
				</button>
			</div>
		</div>
	</div>
	<div class="user-profile" v-if="!loading">
		<!--
		<div class="profile-image">
			<img :src="profile['profile-image-path']" alt="Profile image"/>
		</div>
		-->
		<div class="container-fluid">
			<div class="row">
				<div id="username-box" class="col d-flex align-items-center justify-content-start">
					<div class="user-profile-follow-unfollow" v-if="!banned">
						<button class="btn btn-primary btn-block" @click="followUser" v-if="!followed">Follow</button>
						<button class="btn btn-secondary btn-block" @click="unfollowUser" v-else>Unfollow</button>
					</div>
					<div class="user-profile-ban-unban">
						<button class="btn btn-danger btn-block" @click="banUser" v-if="!banned">Ban</button>
						<button class="btn btn-secondary btn-block" @click="unbanUser" v-else>Unban</button>
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
