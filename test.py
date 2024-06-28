def pack_string(s, encoding='utf-8'):
    """Pack a string into a binary OSC string."""
    if isinstance(s, unicodetype):
        s = s.encode(encoding)
    assert all((i if have_bytes else ord(i)) < 128 for i in s), (
        "OSC strings may only contain ASCII chars.")

    slen = len(s)
    return s + b'\0' * (((slen + 4) & ~0x03) - slen)

def pack_blob(b, encoding='utf-8'):
    """Pack a bytes, bytearray or tuple/list of ints into a binary OSC blob."""
    if isinstance(b, (tuple, list)):
        b = bytearray(b)
    elif isinstance(b, unicodetype):
        b = b.encode(encoding)

    blen = len(b)
    b = pack('>I', blen) + bytes(b)
    return b + b'\0' * (((blen + 3) & ~0x03) - blen)

def create_message(address, *args):
    assert address.startswith('/'), "Address pattern must start with a slash."

    data = []
    types = [',']

    for arg in args:
        type_ = type(arg)

        if isinstance(arg, tuple):
            typetag, arg = arg
        else:
            typetag = TYPE_MAP.get(type_) or TYPE_MAP.get(arg)

        data.append(pack_blob(arg))
        types.append(typetag)

    return pack_string(address) + pack_string(''.join(types)) + b''.join(data)


class Client:
    def __init__(self, host, port=None):
        if port is None:
            if isinstance(host, (list, tuple)):
                host, port = host
            else:
                port = host
                host = '127.0.0.1'

        self.dest = pack_addr((host, port))
        self.sock = None

    def send(self, msg, *args, **kw):
        dest = pack_addr(kw.get('dest', self.dest))

        if not self.sock:
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

        if isinstance(msg, Bundle):
            msg = pack_bundle(msg)
        elif args or isinstance(msg, unicodetype):
            msg = create_message(msg, *args)

        self.sock.sendto(msg, dest)
