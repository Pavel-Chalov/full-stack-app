import { joinClasses } from "../../lib/joinClasses"
import { TypeUiSizes } from "../../types/uiSizes"
import "./button.scss"

interface ButtonProps extends React.DetailedHTMLProps<React.ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement> {
    uiSize: TypeUiSizes
}

export const Button = ({children, className, uiSize, ...props}: ButtonProps) => {
    const classes = joinClasses("button", uiSize, className)

    return (
        <button {...props} className={classes}>{children}</button>
    )
}