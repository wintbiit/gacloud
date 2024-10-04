import {Spin} from "@douyinfe/semi-ui";

function Spinning({ title } : {title: string}) {
    return (
        <Spin size="large"
              tip={title}
              style={{
                    position: "absolute",
                    top: "50%",
                    left: "50%",
                    transform: "translate(-50%, -50%)",
                    width: "800px",
              }}
        />
    )
}

export default Spinning;