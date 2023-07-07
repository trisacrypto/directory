import MemberSelectNetwork from "../MemberNetworkSelect";
import { fireEvent, act, render, screen } from 'utils/test-utils';
import { dynamicActivate } from 'utils/i18nLoaderHelper';

function renderComponent() {
    return render(<MemberSelectNetwork />);
}

describe("MemberSelectNetwork", () => {
    beforeAll(() => {
        act(() => {
          dynamicActivate('en');
        });
      });

    it("should render", () => {
        const { container } = renderComponent();
        expect(container).toMatchSnapshot();
    });

    it("should have mainnet as the default value", () => {
        renderComponent();
        const select = screen.getByTestId("select-network");
        expect(select).toHaveValue("mainnet");
    });

    it("should have testnet as the value when selected", () => {
        renderComponent();
        const select = screen.getByTestId("select-network");
        fireEvent.change(select, { target: { value: "testnet" } });
        expect(select).toHaveValue("testnet");
    }); 
}
);