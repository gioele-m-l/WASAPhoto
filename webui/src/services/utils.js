export async function getImageFile(imagePath) {
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
        this.errormsg = e.toString;
    }
    this.loading = false;
}