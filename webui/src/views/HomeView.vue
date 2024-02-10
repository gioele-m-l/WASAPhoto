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
			page: 0,
        };
    },
    methods: {
		refresh(){
			this.errormsg = null;
			this.page = 0;
			this.photosStream = [];
			this.getMyStream();
		},

		loadMorePhotos(){
			this.page = this.page+1;
			this.getMyStream();
		},

        async getMyStream() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.get("/photos/", {
					params: {
                        	page: this.page, 
                    },
                    headers: {
                        Authorization: this.authToken,
                    }
                });
				if (response.data.photos != null){
					for(let i=0; i<response.data.photos.length; i++){
						let photoID = response.data.photos[i]['photo-id'];
						let ownerID = response.data.photos[i]['owner-id'];
						let ownerUsername = response.data.photos[i]['owner-username'];
						let timestamp = response.data.photos[i]['timestamp'];
						let imagePath = response.data.photos[i]['image-path'];
						let likesCount = response.data.photos[i]['likes-count'];
						let commentsCount = response.data.photos[i]['comments-count'];
						let caption = response.data.photos[i]['caption'];
						let photo = {photoID, ownerID, ownerUsername, timestamp, imagePath, likesCount, commentsCount, caption};
						this.photosStream.push(photo);
					}
				} else {
					if(this.page > 0){
						this.page = this.page - 1;
					}
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
			let ftype = this.$refs.file.files[0].name.split('.')
			ftype = ftype[ftype.length -1]
			if (ftype == "jpeg"){
				ftype = "jpg"
			}
			this.postPhotoFile = new Blob([this.$refs.file.files[0]], { type: "image/"+ftype });
		},

        async uploadPhoto() {
			this.loading = true;
			this.errormsg = null;
			if (this.postPhotoFile == null || (this.postPhotoFile.type != "image/jpg" && this.postPhotoFile.type != "image/png")){
				this.errormsg = "You must select an image file (png/jpg)";
				this.postPhotoFile = null;
				this.loading = false;
				return null;
			}
			if (this.postPhotoFile.size/1024 > 16500){
				// Check if the file size is greater than 16 MB (16500 kB)
				this.errormsg = "Maximum file size: 16 MB";
				this.postPhotoFile = null;
				this.loading = false;
				return null;
			}
			if (this.postPhotoCaption.length > 100) {
				this.errormsg = "The caption must be maximum 100 characters long"
				this.loading = false;
				return null;
			}

			try {
				const formData = new FormData();
				formData.append("caption", this.postPhotoCaption);
				formData.append("image", this.postPhotoFile);
				const headers = { 'Content-Type': 'multipart/form-data', 
									'Authorization': this.authToken,
									'Access-Control-Allow-Origin': '*'}
				let response = await this.$axios.post("/photos/", formData, { headers });
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.postPhotoCaption = "";
			this.postPhotoFile = null;
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
				placeholder="Some text..."
			/>

			<label for="postPhotoFile">File</label>
			<input
				id="postPhotoFile"
				accept=".png, .jpg, .jpeg"
				@change="uploadFile"
				type="file"
				ref="file"
			>
			<button type="submit">Publish</button>
		</form>
	</div>

	<div class="stream-photos">
		<div v-if="photosStream.length>0">
			<div>
				<PhotoCard v-for="photo in photosStream" :key="photo.photoID" :photo="photo" @photoUpdated="refresh"/>
			</div>
			<div  v-if="photosStream.length%20 == 0">
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
</template>

<style>
</style>
