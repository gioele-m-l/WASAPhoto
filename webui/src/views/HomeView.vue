<script>
import PhotoCard from '../components/PhotoCard.vue';

export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            authToken: sessionStorage.getItem("auth-token"),
            userID: sessionStorage.getItem("user-id"),
            photosStream: [],
			postPhotoModalVisibility: false,
			postPhotoCaption: "",
			postPhotoFile: null,
        };
    },
    methods: {
        async getMyStream() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.get("/photos/", {
                    headers: {
                        Authorization: this.authToken,
                    }
                });
                for(let i=0; i<response.data.length; i++){
					let photoID = response.data[i]['photo-id'];
					let ownerID = response.data[i]['owner'];
					let timestamp = response.data[i]['timestamp'];
					let imagePath = response.data[i]['image-path'];
					let likesCount = response.data[i]['likes-count'];
					let commentsCount = response.data[i]['comments-count'];
					let caption = response.data[i]['caption'];
					let photo = {photoID, ownerID, timestamp, imagePath, likesCount, commentsCount, caption};
					this.photosStream.push(photo);
				}
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },

		showPostPhotoModal(){
			this.postPhotoModalVisibility = true;
		},

		hidePostPhotoModal(){
			this.postPhotoModalVisibility = false;
		},

		uploadFile(){
			this.postPhotoFile = this.$refs.file.files[0];
		},

        async uploadPhoto() {
			this.loading = true;
			this.errormsg = false;

			try {
				const formData = new FormData();
				formData.append("caption", this.postPhotoCaption);
				formData.append("image", this.postPhotoFile);
				let response = await this.$axios.post("/photos/", formData, { 
					headers : {
						Authorization: this.authToken,
						Content: 'multipart/form-data'
					}
				});

				console.log(response.status)
			} catch (e) {
				this.errormsg = e.toString();
			}

			this.hidePostPhotoModal();
			this.loading = false;
        },
    },
    mounted() {
        this.getMyStream();
    },
    components: { PhotoCard }
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Stream</h1>
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
					<button type="button" class="btn btn-sm btn-outline-primary" @click="showPostPhotoModal">
						Post a photo
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>

	<div class="post-photo-box" v-if="postPhotoModalVisibility">
		<h6>Choose a file</h6>
		<button @click="hidePostPhotoModal">&times;</button>
		<form @submit.prevent="uploadPhoto">
			<label for="postPhotoCaption">Caption</label>
			<input
				id="postPhotoCaption"
				v-model="postPhotoCaption"
				type="text"
			/>

			<label for="postPhotoFile">File</label>
			<input
				id="postPhotoFile"
				@change="uploadFile"
				type="file"
				ref="file"
			>
			<button type="submit">Publish</button>
		</form>
	</div>

	<div class="stream-photos">
		<PhotoCard v-for="photo in photosStream" :key="photo.photoID" :photo="photo" v-if="photosStream.length != 0" />
		<h5 v-else>There are no photos yet :'(</h5>
	</div>
</template>

<style>
</style>
