import $ from "jquery"

export const animateForm = () => {
    $('.auth-form').animate({height: "toggle", opacity: "toggle"}, "slow");
}