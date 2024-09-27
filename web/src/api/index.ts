import axios from "axios";

const axioser = axios.create({
    baseURL: "/api/v1",
    timeout: 5000,
})

axioser.interceptors.request.use((config) => {
    // getAuthToken from localStorage
    const authToken = localStorage.getItem("authToken")
    if (authToken) {
        config.headers.Authorization = `Bearer ${authToken}`
    }

    return config
})

axioser.interceptors.response.use((response) => {
    return response
}, (error) => {
    console.log("Error response", error.response)

    if (error.response.status === 401 && !window.location.pathname.startsWith("/login")) {
        console.log("Unauthorized, redirecting to login")

        localStorage.removeItem("authToken")
        const location = window.location.pathname

        document.location.href = "/login?redirect=" + encodeURIComponent(location)
    }

    if (error.response.status === 503 && !window.location.pathname.startsWith("/maintenance")) {
        document.location.href = "/maintenance"
    }

    if (error.response.status === 530 && !window.location.pathname.startsWith("/setup")) {
        document.location.href = "/setup"
    }

    return Promise.reject(error)
})

const logout = () => {
    localStorage.removeItem("authToken")
    document.location.href = "/login"

    return {
        pathname: "/login"
    }
}

export { axioser, logout }