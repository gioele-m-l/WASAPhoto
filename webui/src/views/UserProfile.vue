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
            } catch(e){
                this.errormsg = e.toString();
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
			<img :src="profile['profile-image-path']" alt="Profile image"/>
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
            <div class="user-profile-follow-unfollow">
                <button @click="followUser">Follow</button>
                <button @click="unfollowUser">Unfollow</button>
            </div>
            <div class="user-profile-ban-unban">
                <button @click="banUser">Ban</button>
                <button @click="unbanUser">Unban</button>
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
