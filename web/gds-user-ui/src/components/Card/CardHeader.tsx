import { Heading, HeadingProps } from '@chakra-ui/react';

const CardHeader: React.FC<HeadingProps> = ({ children, ...props }) => {
  return (
    <Heading size="sm" {...props}>
      {children}
    </Heading>
  );
};

export default CardHeader;
