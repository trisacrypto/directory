import { Meta, Story } from "@storybook/react";
import MobileNav from "./MobileNav";

interface MobileProps {
  onOpen: () => void;
}

export default {
  title: "components/MobileNav",
  component: MobileNav,
} as Meta<MobileProps>;

const Template: Story<MobileProps> = (args) => <MobileNav {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
