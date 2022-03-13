import { Box, BoxProps } from "@chakra-ui/react";

interface CardProps extends BoxProps {
  children: React.ReactNode;
}
const Card: React.FC<CardProps> = (props) => {
  const { children, ...rest } = props;
  return (
    <Box
      border="2px solid #E5EDF1"
      borderRadius="10px"
      padding={{ base: 3, md: 9 }}
      fontFamily="Open Sans"
      {...rest}
    >
      {children}
    </Box>
  );
};

export default Card;
