using System;
using System.Net.Sockets;
using System.Text;
using Microsoft.Xna.Framework;
using Microsoft.Xna.Framework.Graphics;
using Microsoft.Xna.Framework.Input;

namespace MazeGame;

internal enum ServerOpCodes {
    CmdConnect = 0,
    CmdDisconnect = 1,
    CmdMove = 2,
    CmdRequestPosition = 3,
}

public class GameManager : Game {
    private GraphicsDeviceManager _graphics;
    private SpriteBatch _spriteBatch;

    public GameManager() {
        _graphics = new GraphicsDeviceManager(this);
        Content.RootDirectory = "Content";
        IsMouseVisible = true;
    }

    protected override void Initialize() {
        // TODO: Add your initialization logic here

        // Test server connection by sending a dummy packet
        byte[] opcode = BitConverter.GetBytes((int)ServerOpCodes.CmdConnect);

        byte[] jsonBytes
            = Encoding.UTF8.GetBytes(
                @"{""name"": ""Jake"", ""color"": 0, ""roomName"": ""Jake's Room""}");
        byte[] dataToSend = new byte[opcode.Length + jsonBytes.Length];
        Buffer.BlockCopy(opcode, 0, dataToSend, 0, opcode.Length);
        Buffer.BlockCopy(jsonBytes, 0, dataToSend, opcode.Length, jsonBytes.Length);

        using (TcpClient client = new TcpClient("127.0.0.1", 40200))
        using (NetworkStream stream = client.GetStream()) {
            stream.Write(dataToSend, 0, dataToSend.Length);
            Console.WriteLine("Sending data to server");

            // Optionally, read a response
            // byte[] buffer = new byte[1024];
            // int bytesRead = stream.Read(buffer, 0, buffer.Length);
            // Console.WriteLine("Response from server: " + Encoding.UTF8.GetString(buffer, 0, bytesRead));
        }

        base.Initialize();
    }

    protected override void OnExiting(object sender, ExitingEventArgs e) {
        byte[] opcode = BitConverter.GetBytes((int)ServerOpCodes.CmdDisconnect);

        using (TcpClient client = new TcpClient("127.0.0.1", 40200))
        using (NetworkStream stream = client.GetStream()) {
            stream.Write(opcode, 0, opcode.Length);
            Console.WriteLine("Sending data to server");
        }

        base.OnExiting(sender, e);
    }

    protected override void LoadContent() {
        _spriteBatch = new SpriteBatch(GraphicsDevice);

        // TODO: use this.Content to load your game content here
    }

    protected override void Update(GameTime gameTime) {
        if (GamePad.GetState(PlayerIndex.One).Buttons.Back == ButtonState.Pressed ||
            Keyboard.GetState().IsKeyDown(Keys.Escape))
            Exit();

        // TODO: Add your update logic here

        base.Update(gameTime);
    }

    protected override void Draw(GameTime gameTime) {
        GraphicsDevice.Clear(Color.CornflowerBlue);

        // TODO: Add your drawing code here

        base.Draw(gameTime);
    }
}