<script>
import PhotoCard from '../components/PhotoCard.vue';
import UserCard from '../components/UserCard.vue';

export default {
	components: {
		PhotoCard,
		UserCard,
	},
	
	watch: {
    	'$route.params.username': function(newUsername, oldUsername) {
			this.username = newUsername;
      		this.refresh();
		}
    },

	data: function() {
		return {
			errormsg: null,
			loading: false,
            sessionUsername: sessionStorage.getItem("username"),
			sessionUserID: sessionStorage.getItem("user-id"),
			sessionAuthToken: sessionStorage.getItem("auth-token"),

            username: this.$route.params.username,
			profile: {},
			photos: [],
			profileImage: null,
			page: 0,

			followed: false,
			banned: false,
			found: false,

			showFollowingVar: false,
			showFollowersVar: false,

			followers: [],
			following: [],
		}
	},
	methods: {
		logout(){
			sessionStorage.clear();
			this.$router.replace({ path: "/login"});
		},

		refresh(){
			this.showFollowingVar = false;
			this.showFollowersVar = false;
			this.banned = false;
			this.followed = false;
			this.profile = {};
			this.photos = [];
			this.page = 0;
			this.followers = [];
			this.following = [];
			this.errormsg = null;

			this.getUserProfile();
			this.checkBan();
			if(this.banned != true){
				this.checkFollow();
			}
			this.listFollowers();
			this.listFollowing();
			this.getUserPhotos();
		},

		async checkBan(){
			this.loading = true;
			this.errormsg = null;

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
			this.errormsg = null;

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

		async getUserProfile() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.username + "/", {
						headers: {
							Authorization: this.sessionAuthToken,
						}
					} 
				);
				
				this.profile = response.data;
				this.found = true;
				this.getImageFile();
			} catch (e) {
				if (e.response.status == 404){
					this.found = false;
				} else {
					this.errormsg = e.toString();
				}
				
			}
			this.loading = false;	
		},

		loadMorePhotos(){
			this.page = this.page+1;
			this.getUserPhotos();
		},

		async getUserPhotos(){
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + this.username + "/photos/", {
						params: {
                        	page: this.page, 
                    	},
						headers: {
							Authorization: this.sessionAuthToken,
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
                if(e.response.status != 404){
				    this.errormsg = e.toString();
                }
			}
			this.loading = false;
		},

        async followUser(){
            this.loading = true;
            this.errormsg = null;
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
            this.errormsg = null;
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
            this.errormsg = null;
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
            this.errormsg = null;
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

		async getImageFile(){
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/images/" + this.profile['profile-image-path'], {
						headers: {
							Authorization: this.sessionAuthToken,
						}
					}
				);
				let ext = response.headers['content-type'].split('/')[1];
				this.profileImage = 'data:image/'+ext+';base64,'+response.data;
			} catch(e) {
				if (e.response.status != 404){
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		},

		async listFollowing(){
			this.loading = true;
			this.errormsg = null;

			try {
				let response = await this.$axios.get("/users/" + this.username + "/followings/", {
					headers : {
						Authorization: this.sessionAuthToken,
					}
				})

				if(response.data != null){
					for(let i=0; i<response.data.length; i++){
						let user = {
							userID: response.data[i]['user-id'],
							username: response.data[i]['username'],
							imagePath: response.data[i]['profile-image-path'],
						};
						this.following.push(user);
					}
				}
			} catch (e){
				this.errormsg = e.toString();
			}

			this.loading = false;
		},

		async listFollowers(){
			this.loading = true;
			this.errormsg = null;

			try {
				let response = await this.$axios.get("/users/" + this.username + "/followers/", {
					headers : {
						Authorization: this.sessionAuthToken,
					}
				})

				if(response.data != null){
					for(let i=0; i<response.data.length; i++){
						let user = {
							userID: response.data[i]['user-id'],
							username: response.data[i]['username'],
							imagePath: response.data[i]['profile-image-path'],
						};
						this.followers.push(user);
					}
				}
			} catch (e){
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		showFollowers(){
			this.showFollowingVar = false;
			this.showFollowersVar = true;
		},

		showFollowing(){
			this.showFollowingVar = true;
			this.showFollowersVar = false;
		},

		hideSideBar(){
			this.showFollowingVar = false;
			this.showFollowersVar = false;
		},
		
	},

	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div v-if="found">
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">{{ profile['username'] }}'s profile</h1>
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
		<div class="row vh-100 w-100">
			<div class="col">
				<div class="user-profile">
					<div class="container-fluid">
						<div class="row">
							<div class="col">
								<div v-if="profileImage != null">
									<img :src="profileImage" alt="Profile image" style="width: 130px; height: 125px; border-radius: 50%; object-fit: cover; border: 1px solid #000;"/>
								</div>
								<div v-else>
									<img src="/default-user.jpg" style="width: 130px; height: 125px; border-radius: 50%; object-fit: cover; border: 1px solid #000;">
								</div>
							</div>
							<div id="user-stats" class="col">
								<ul class="list-group list-group-horizontal">
									<li class="list-group-item text-center w-50">
										Photos
										<br>
										{{ profile['photos-count'] }}
									</li>
									<li class="list-group-item text-center w-50"  @click="showFollowers">
										Followers
										<br>
										{{ followers.length }}
									</li>
									<li class="list-group-item text-center w-50"  @click="showFollowing">
										Following
										<br>
										{{ following.length }}
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
								<div  v-if="!banned" class="user-profile-follow-unfollow">
									<div v-if="!followed">
										<button class="btn btn-primary btn-block" @click="followUser">Follow</button>
									</div>
									<div v-else>
										<button class="btn btn-secondary btn-block" @click="unfollowUser">Unfollow</button>
									</div>
								</div>
								<div class="user-profile-ban-unban">
									<div v-if="!banned">
										<button class="btn btn-danger btn-block" @click="banUser">Ban</button>
									</div>
									<div v-else>
										<button class="btn btn-secondary btn-block" @click="unbanUser">Unban</button>
									</div>
								</div>
							</div>
						</div>
					</div>
					<hr>
					<div class="profile-photos">
						<h3>Photos</h3>
						<div v-if="photos.length>0">
							<div>
								<PhotoCard v-for="photo in photos" :key="photo.photoID" :photo="photo" @photoUpdated="refresh"/>
							</div>
							<div  v-if="profile['photos-count']%20 != 0 && Math.floor(profile['photos-count']/20) != page">
								<button class="btn btn-outline-secondary" @click="loadMorePhotos">More photos</button>
							</div>
							<div v-else>
								<p>-- End of Feed --</p>
							</div>
						</div>
						<div v-else>
							<h5>There are no photos yet :'(</h5>
							<!--
								<h6>Post something<img src="https://i.redd.it/4s978dxj7xp51.jpg" style="width: 100px; heigth: 100px;"></h6>
							-->
						</div>
					</div>
				</div>
			</div>
			<div  v-if="showFollowersVar || showFollowingVar" class="col-4 border">
				<div v-if="showFollowersVar">
					<h3 class="d-flex border-bottom justify-content-between">Followers<button @click="hideSideBar" class="btn btn-icon btn-sm">&times;</button></h3>
					<div v-if="followers.length > 0" class="w-75">
						<UserCard v-for="user in followers" :key="user.userID" :user="user"></UserCard>
					</div>
				</div>
				<div v-else-if="showFollowingVar">
					<h3 class="d-flex border-bottom justify-content-between">Following<button @click="hideSideBar" class="btn btn-icon btn-sm">&times;</button></h3>
					<div v-if="following.length > 0" class="w-75">
						<UserCard v-for="user in following" :key="user.userID" :user="user"></UserCard>
					</div>
				</div>
			</div>
		</div>
	</div>
		
	<div v-else-if="!loading">
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Error 404: user "{{ username }}" was not found...</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>
	</div>
	
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
</style>
