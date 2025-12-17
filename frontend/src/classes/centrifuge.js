import { Centrifuge, UnauthorizedError } from 'centrifuge';

export class OPCentrifuge {
  centrifuge = null;

  async init(token = null) {
    if (this.centrifuge) {
      // await this.centrifuge.ready(2500);
      return;
    }

    let protocol = 'ws';
    let url = import.meta.env.VITE_CENTRIFUGO_WS;

    if (window.location.protocol === 'https:') {
      protocol = 'wss';
      url = import.meta.env.VITE_CENTRIFUGO_WSS;
    }

    this.centrifuge = new Centrifuge(`${protocol}://${url}/connection/websocket`, {
      getToken: this.generateConnectionToken
    });

    if (token) {
      this.centrifuge.setToken(token);
    }

    this.centrifuge.connect();
  }

  async generateConnectionToken() {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/broadcasting/connect`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      if (!res.ok) {
        if (res.status === 403) {
          // Return special error to not proceed with token refreshes, client will be disconnected.
          throw new UnauthorizedError();
        }
        // Any other error thrown will result into token refresh re-attempts.
        throw new Error(`Unexpected status code ${res.status}`);
      }
      const data = await res.json();
      return data.token;
    } catch (e) {
      console.log(e);
    }
  }

  async getSubscriptionToken(ctx) {
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/broadcasting/auth`, {
        method: "POST",
        body: JSON.stringify(ctx),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
        },
      });
      const data = await res.json();
      return data.token;
    } catch (e) {
      if (e.response && e.response.status === 403) {
        console.error(e)
        return "";
      }

      throw e;
    }
  }

  newSubscription (name, options) {
    const subscriptionOptions = {...options};

    // Inget behov av subscription token för user channels, men bör kunna hanteras snyggare?
    if (name.indexOf('#') === -1) {
      subscriptionOptions.getToken = (ctx) => this.getSubscriptionToken(ctx);
    }

    return this.centrifuge.newSubscription(name, {
      ...subscriptionOptions,
    });
  }

  removeSubscription(sub) {
    if (sub) {
      this.centrifuge.removeSubscription(sub);
    }
  }
}

export const centrifuge = new OPCentrifuge();

