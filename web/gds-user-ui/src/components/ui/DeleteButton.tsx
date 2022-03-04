import { DeleteIcon } from "@chakra-ui/icons";
import { Button, Icon, Tooltip, TooltipProps } from "@chakra-ui/react";

type DeleteButtonProps = {
  tooltip?: Omit<TooltipProps, "children">;
};

const TOOLTIPS_DELAY = 2 * 1000;

const DeleteButton: React.FC<DeleteButtonProps> = ({ tooltip }) => {
  return (
    <Tooltip label="Delete" openDelay={TOOLTIPS_DELAY} {...tooltip}>
      <Button
        variant="ghost"
        _hover={{ background: "red.100", color: "red.500" }}
        _focus={{ background: "red.100", color: "red.500" }}
        borderRadius={0}
      >
        <Icon as={DeleteIcon} />
      </Button>
    </Tooltip>
  );
};

export default DeleteButton;
