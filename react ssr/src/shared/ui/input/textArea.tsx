import { joinClasses } from "../../lib/joinClasses";
import { TypeUiSizes } from "../../types/uiSizes";
import "./input.scss"
import {ForwardedRef, forwardRef, useState} from "react"

type TextAreaProps = {
    uiSize: TypeUiSizes;
    error?: string;
    label?: string
} & React.DetailedHTMLProps<React.TextareaHTMLAttributes<HTMLTextAreaElement>, HTMLTextAreaElement>

export const TextArea = forwardRef<HTMLTextAreaElement, TextAreaProps>((
  {
    className,
    error,
    label,
    uiSize,
    placeholder,
    value,
    ...props
  },
  ref: ForwardedRef<HTMLTextAreaElement>
) => {

  const [hasValue, setHasValue] = useState(value?.toString().trim() !== "" && value !== undefined)

  const classes = joinClasses("field", error ? "error" : undefined, uiSize, className);

  return (
    <div className={classes}>
      <textarea
        ref={ref}
        onChange={e => setHasValue(e.target.value.trim() !== "")}
        autoComplete="off"
        {...props}
        value={value}
        className={joinClasses("field__input", hasValue ? "hasValue" : "")}
        placeholder={placeholder}
      />
      <div className="field__bar"></div>
      <label className="field__label" htmlFor="name">{label}</label>
      <small className="field__hint">{error}</small>
    </div>
  );
});