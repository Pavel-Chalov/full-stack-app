import { joinClasses } from "../../lib/joinClasses";
import { TypeUiSizes } from "../../types/uiSizes";
import "./input.scss"
import {ForwardedRef, forwardRef} from "react"

type InputProps = {
    uiSize: TypeUiSizes;
    error?: string;
    label?: string
} & React.DetailedHTMLProps<React.InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>

export const Input = forwardRef<HTMLInputElement, InputProps>((
    {
      className,
      error,
      label,
      uiSize,
      placeholder,
      value,
      ...props
    },
    ref: ForwardedRef<HTMLInputElement>
  ) => {

    const classes = ["form-input", error ? "form-has-error" : null, uiSize, className].join(" ");

  return (
    <div className={classes}>
      <input
        ref={ref}
        autoComplete="off"
        {...props}
        value={value}
        className={joinClasses("form-element-input", value?.toString().trim() !== "" ? "hasValue" : "")}
        placeholder={placeholder}
      />
      <div className="form-element-bar"></div>
      <label className="form-element-label">{label}</label>
      <small className="form-element-hint">{error}</small>
    </div>
  );
});