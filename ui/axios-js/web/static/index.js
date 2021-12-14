const initializeAxiosInstance = async (url) => {
    var csrfToken;
    try {
        let resp = await axios.get(url, {withCredentials: true});
        console.log(resp);
        csrfToken = resp.headers["x-csrf-token"];
        console.log(csrfToken);
        return axios.create({
            withCredentials: true,
            headers: {"X-CSRF-Token": csrfToken}
        });
    } catch (err) {
        console.log(err);
    }
};

const post = async (axiosInstance, url) => {
    try {
        let resp = await axiosInstance.post(url);
        console.log(resp);
    } catch (err) {
        console.log(err);
    }
};

const url = "http://localhost:8080/api";
initializeAxiosInstance(url)
    .then(axiosInstance => {
        post(axiosInstance, url);
    });
