import React from 'react';
import "./title.scss";
import { joinClasses } from '../../lib/joinClasses';

interface TitleProps<T extends keyof JSX.IntrinsicElements = 'h1' | "h2" | "h3" | "h4" | "h5" | "h6"> extends React.HTMLAttributes<HTMLElement> {
  Level: T;
  className?: string;
}

export const Title: React.FC<TitleProps> = ({ 
  className, 
  children, 
  Level = 'h1', 
  ...props 
}) => {
  const classes = joinClasses("title", className)

  return (
    <Level className={classes} {...props}>
      {children}
    </Level>
  );
};
