import * as React from 'react';
import ContentLoader from 'react-content-loader';

const Avatar = (props: any) => (
  <ContentLoader
    viewBox="0 0 660 475"
    // height={70}
    // width={400}
    // speed={2}
    primaryColor="#f3f3f3"
    secondaryColor="#ecebeb"
    {...props}>
    <circle cx="420" cy="30" r="16" />
    <rect x="379" y="18" rx="5" ry="0" width="20" height="10" />
    <rect x="379" y="35" rx="5" ry="0" width="20" height="10" />
  </ContentLoader>
);

export default Avatar;
