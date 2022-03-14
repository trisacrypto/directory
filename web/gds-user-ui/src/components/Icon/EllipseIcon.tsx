import * as React from "react";

interface SvgTrisaIconProps extends React.SVGProps<SVGEllipseElement> {
  fillCurrent?: boolean | undefined;
}

const SvgEllipseIcon = (props: SvgTrisaIconProps) => {
  return (
    <svg width="28" height="30" viewBox="0 0 28 30" fill="currentColor">
      <ellipse
        rx="13.4735"
        ry="15"
        transform="matrix(-1 0 0 1 13.6183 15)"
        fill="#60C4CA"
        {...props}
      />
    </svg>
  );
};

export default SvgEllipseIcon;
